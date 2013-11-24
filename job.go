package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
	"github.com/wkharold/jobd/deps/github.com/gorhill/cronexpr"

	"fmt"
	"regexp"
	"strings"
	"time"
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
	state    string
}

type jobreader func() []byte
type jobwriter func([]byte) (int, error)

type job struct {
	srv.File
	defn jobdef
	done chan bool
}

type jobfile struct {
	srv.File
	reader jobreader
	writer jobwriter
}

func mkJob(root *srv.File, user p.User, def jobdef) (*job, error) {
	glog.V(4).Infof("Entering mkJob(%v, %v, %v)", root, user, def)
	defer glog.V(4).Infof("Exiting mkJob(%v, %v, %v)", root, user, def)

	glog.V(3).Infoln("Creating job directory: ", def.name)

	job := &job{defn: def, done: make(chan bool)}

	ctl := &jobfile{
		reader: func() []byte {
			return []byte(job.defn.state)
		},
		writer: func(data []byte) (int, error) {
			switch cmd := strings.ToLower(string(data)); cmd {
			case STOP:
				if job.defn.state != STOPPED {
					glog.V(3).Infof("Stopping job: %v", job.defn.name)
					job.defn.state = STOPPED
					job.done <- true
				}
				return len(data), nil
			case START:
				if job.defn.state != STARTED {
					glog.V(3).Infof("Starting job: %v", job.defn.name)
					job.defn.state = STARTED
					go job.run()
				}
				return len(data), nil
			default:
				return 0, fmt.Errorf("Unknown command: %s", cmd)
			}
		}}
	if err := ctl.Add(&job.File, "ctl", user, nil, 0666, ctl); err != nil {
		glog.Errorf("Can't create %s/ctl [%v]", def.name, err)
		return nil, err
	}

	sched := &jobfile{
		reader: func() []byte {
			return []byte(def.schedule)
		},
		writer: func(data []byte) (int, error) {
			return 0, srv.Eperm
		}}
	if err := sched.Add(&job.File, "schedule", user, nil, 0444, sched); err != nil {
		glog.Errorf("Can't create %s/schedule [%v]", def.name, err)
		return nil, err
	}

	cmd := &jobfile{
		reader: func() []byte {
			return []byte(def.cmd)
		},
		writer: func(data []byte) (int, error) {
			return 0, srv.Eperm
		}}
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

	return &jobdef{name, schedule, cmd, STOPPED}, nil
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

func (jf jobfile) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jofile.Read(%v, %v, %)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobfile.Read(%v, %v, %v)", fid, buf, offset)

	cont := jf.reader()

	if offset > uint64(len(cont)) {
		return 0, nil
	}

	contout := cont[offset:]

	copy(buf, contout)
	return len(contout), nil
}

func (jf jobfile) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering jobfile.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting jobfile.Wstat(%v, %v, %v)", fid, dir)

	return nil
}

func (jf *jobfile) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobctl.Write(%v, %v, %v)", fid, data, offset)
	defer glog.V(4).Infof("Exiting jobctl.Write(%v, %v, %v)", fid, data, offset)

	jf.Parent.Lock()
	defer jf.Parent.Unlock()

	return jf.writer(data)
}

func (j *job) run() {
	for {
		now := time.Now()
		e, _ := cronexpr.Parse(j.defn.schedule)
		select {
		case <-time.After(e.Next(now).Sub(now)):
			glog.V(3).Infof("running `%s`", j.defn.cmd)
		case <-j.done:
			glog.V(3).Infof("completed")
			return
		}
	}
}
