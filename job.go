package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
)

type job struct {
	srv.File
}

func mkJob(root *srv.File, user p.User, jobname string) (*job, error) {
	glog.V(4).Infof("Entering mkJob(%v, %v, %v)", root, user, jobname)
	defer glog.V(4).Infof("Exiting mkJob(%v, %v, %v)", root, user, jobname)

	glog.V(3).Infoln("Creating job directory: ", jobname)

	job := &job{}
	if err := job.Add(root, jobname, user, nil, p.DMDIR|0555, job); err != nil {
		glog.Errorf("Can't create job %s: %v\n", jobname, err)
		return nil, err
	}

	return job, nil
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
