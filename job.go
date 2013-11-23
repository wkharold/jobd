package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
	"github.com/wkharold/jobd/deps/github.com/gorhill/cronexpr"

	"fmt"
	"regexp"
	"strings"
)

const (
	STOPPED = "stopped"
	STOP    = "stop"
	STARTED = "started"
	START   = "start"
)

type jobdef struct {
	name     string
	schedule string
	cmd      string
}

type job struct {
	srv.File
}

type jobctl struct {
	srv.File
	state string
}

type jobsched struct {
	srv.File
	schedule string
}

type jobcmd struct {
	srv.File
	cmd string
}

func mkJob(root *srv.File, user p.User, def jobdef) (*job, error) {
	glog.V(4).Infof("Entering mkJob(%v, %v, %v)", root, user, def)
	defer glog.V(4).Infof("Exiting mkJob(%v, %v, %v)", root, user, def)

	glog.V(3).Infoln("Creating job directory: ", def.name)

	job := &job{}

	ctl := &jobctl{state: STOPPED}
	if err := ctl.Add(&job.File, "ctl", user, nil, 0555, ctl); err != nil {
		glog.Errorf("Can't create %s/ctl [%v]", def.name, err)
		return nil, err
	}

	sched := &jobsched{schedule: def.schedule}
	if err := sched.Add(&job.File, "schedule", user, nil, 0444, sched); err != nil {
		glog.Errorf("Can't create %s/schedule [%v]", def.name, err)
		return nil, err
	}

	cmd := &jobcmd{cmd: def.cmd}
	if err := cmd.Add(&job.File, "cmd", user, nil, 0444, cmd); err != nil {
		glog.Errorf("Can't create %s/cmd [%v]", def.name, err)
		return nil, err
	}

	return job, nil
}

func mkJobDefinition(name, schedule, cmd string) (*jobdef, error) {
	if ok, err := regexp.MatchString("[^[:word:]]", name); ok || err != nil {
		switch {
		case ok:
			return nil, fmt.Errorf("Invalid job name: %s", name)
		default:
			return nil, err
		}
	}

	if _, err := cronexpr.Parse(schedule); err != nil {
		return nil, err
	}

	return &jobdef{name, schedule, cmd}, nil
}

func (j job) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering job.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Entering job.Read(%v, %v, %v)", fid, buf, offset)

	return 0, nil
}

func (j *job) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering job.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting job.Wstat(%v, %v)", fid, dir)

	return nil
}

func (ctl jobctl) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobctl.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobctl.Read(%v, %v, %v)", fid, buf, offset)

	if offset > uint64(len(ctl.state)) {
		return 0, nil
	}

	copy(buf, ctl.state[offset:])
	return len(ctl.state[offset:]), nil
}

func (ctl *jobctl) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobctl.Write(%v, %v, %v)", fid, data, offset)
	defer glog.V(4).Infof("Exiting jobctl.Write(%v, %v, %v)", fid, data, offset)

	ctl.Lock()
	defer ctl.Lock()

	switch cmd := strings.ToLower(string(data)); cmd {
	case STOP:
		if ctl.state != STOPPED {
			glog.V(3).Infof("Stopping job: %v", ctl.Parent)
		}
		return len(data), nil
	case START:
		if ctl.state != STARTED {
			glog.V(3).Infof("Starting job: %v", ctl.Parent)
		}
		return len(data), nil
	default:
		return 0, fmt.Errorf("Unknown command: %s", cmd)
	}
}

func (ctl jobctl) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering jobctl.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting jobctl.Wstat(%v, %v)", fid, dir)

	return nil
}

func (sched jobsched) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobsched.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobsched.Read(%v, %v, %v)", fid, buf, offset)

	if offset > uint64(len(sched.schedule)) {
		return 0, nil
	}

	copy(buf, sched.schedule[offset:])
	return len(sched.schedule[offset:]), nil
}

func (sched jobsched) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering jobsched.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting jobsched.Wstat(%v, %v)", fid, dir)

	return nil
}

func (cmd jobcmd) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobcmd.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobcmd.Read(%v, %v, %v)", fid, buf, offset)

	if offset > uint64(len(cmd.cmd)) {
		return 0, nil
	}

	copy(buf, cmd.cmd[offset:])
	return len(cmd.cmd[offset:]), nil
}

func (cmd *jobcmd) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering jobcmd.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting jobcmd.Wstat(%v, %v)", fid, dir)

	return nil
}
