package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
)

type clonefile struct {
	srv.File
}

func mkCloneFile(dir *srv.File, user p.User) error {
	glog.V(4).Infoln("Entering mkCloneFile(%v, %v)", dir, user)
	defer glog.V(4).Infoln("Exiting mkCloneFile(%v, %v)", dir, user)

	glog.V(3).Infoln("Create the clone file")

	k := new(clonefile)
	if err := k.Add(dir, "clone", user, nil, 0666, k); err != nil {
		glog.Errorln("Can't create clone file: ", err)
		return err
	}

	return nil
}

func (k *clonefile) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	glog.V(4).Infof("Entering clonefile.Write(%v, %v, %v)", fid, data, offset)
	defer glog.V(4).Infof("Exiting clonefile.Write(%v, %v, %v)", fid, data, offset)

	k.Lock()
	defer k.Unlock()

	jobdef := string(data)
	glog.V(3).Infof("Create a new job from: %s", jobdef)

	if err := jobsroot.addJob(jobdef); err != nil {
		return len(data), err
	}

	return len(data), nil
}

func (k *clonefile) Wstat(fid *srv.FFid, dir *p.Dir) error {
	glog.V(4).Infof("Entering clonefile.Wstat(%v, %v)", fid, dir)
	defer glog.V(4).Infof("Exiting clonefile.Wstat(%v, %v)", fid, dir)

	return nil
}
