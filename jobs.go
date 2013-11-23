package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
)

type jobsdir struct {
	srv.File
	user p.User
}

func mkJobsDir(dir *srv.File, user p.User) (*jobsdir, error) {
	glog.V(4).Infof("Entering mkJobsDir(%v, %v)", dir, user)
	defer glog.V(4).Infof("Leaving mkJobsDir(%v, %v)", dir, user)

	glog.V(3).Infoln("Create the jobs directory")

	jobs := &jobsdir{user: user}
	if err := jobs.Add(dir, "jobs", user, nil, p.DMDIR|0555, jobs); err != nil {
		glog.Errorln("Can't create jobs directory ", err)
		return nil, err
	}

	return jobs, nil
}

func (jd jobsdir) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering jobsdir.Read(%v, %v, %v)", fid, buf, offset)
	defer glog.V(4).Infof("Exiting jobsdir.Read(%v, %v, %v)", fid, buf, offset)

	return 0, nil
}

func (jd *jobsdir) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering jobsdir.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting jobsdir.Wstat(%v, %v)", fid, dir)
	return nil
}

func (jd *jobsdir) addJob(def jobdef) error {
	glog.V(4).Infof("Entering jobsdir.addJob(%s)", def)
	defer glog.V(4).Infof("Leaving jobsdir.addJob(%s)", def)

	glog.V(3).Info("Add job: ", def)

	job, err := mkJob(&jd.File, jd.user, def)
	if err != nil {
		return err
	}

	if err := job.Add(&jd.File, def.name, jd.user, nil, p.DMDIR|0444, job); err != nil {
		glog.Errorf("Can't add job %s to jobs directory", def.name)
		return err
	}

	return nil
}
