package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
)

const (
	STOPPED = "stopped"
	STARTED = "started"
)

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

func mkJob(root *srv.File, user p.User, name, schedule, command string) error {
	glog.V(4).Infof("Entering mkJob(%v, %v, %s, %s, %s)", root, user, name, schedule, command)
	defer glog.V(4).Infof("Exiting mkJob(%v, %v, %s, %s, %s)", root, user, name, schedule, command)

	glog.V(3).Infoln("Creating job directory: ", name)

	job := &job{}
	if err := job.Add(root, name, user, nil, p.DMDIR|0444, job); err != nil {
		glog.Errorf("Can't add job directory %s to jobs", name)
		return err
	}

	ctl := &jobctl{state: STOPPED}
	if err := ctl.Add(&job.File, "ctl", user, nil, 0555, ctl); err != nil {
		glog.Errorf("Can't create %s/ctl [%v]", name, err)
		return err
	}

	sched := &jobsched{schedule: schedule}
	if err := sched.Add(&job.File, "schedule", user, nil, 0444, sched); err != nil {
		glog.Errorf("Can't create %s/schedule [%v]", name, err)
		return err
	}

	cmd := &jobcmd{cmd: command}
	if err := cmd.Add(&job.File, "cmd", user, nil, 0444, cmd); err != nil {
		glog.Errorf("Can't create %s/cmd [%v]", name, err)
		return err
	}

	return nil
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

func (sched jobsched) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobsched.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobsched.Read(%v, %v, %v)", fid, buf, offset)

	if offset > uint64(len(sched.schedule)) {
		return 0, nil
	}

	copy(buf, sched.schedule[offset:])
	return len(sched.schedule[offset:]), nil
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
