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

// mkJobsDir create the jobs directory at the root of the jobd name space.
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

// addJob uses mkJob to create a new job subtree for the given job definition and adds it to
// the jobd name space under the jobs directory.
func (jd *jobsdir) addJob(def jobdef) error {
	glog.V(4).Infof("Entering jobsdir.addJob(%s)", def)
	defer glog.V(4).Infof("Leaving jobsdir.addJob(%s)", def)

	glog.V(3).Info("Add job: ", def)

	job, err := mkJob(&jd.File, jd.user, def)
	if err != nil {
		return err
	}

	if err := job.Add(&jd.File, def.name, jd.user, nil, p.DMDIR|0555, job); err != nil {
		glog.Errorf("Can't add job %s to jobs directory", def.name)
		return err
	}

	return nil
}
