package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	_ "github.com/wkharold/jobd/deps/github.com/golang/glog"

	"flag"
	"os"
)

// jobsroot is the root of the jobd file hierarchy
var jobsroot *jobsdir

func main() {
	flfsaddr := flag.String("fsaddr", "0.0.0.0:5640", "Address where job file service listens for connections")
	fldebug := flag.Bool("debug", false, "9p debugging to stderr")
	flag.Parse()

	root, err := mkjobfs()
	if err != nil {
		os.Exit(1)
	}

	s := srv.NewFileSrv(root)
	s.Dotu = true
	if *fldebug {
		s.Debuglevel = 1
	}
	s.Start(s)

	if err := s.StartNetListener("tcp", *flfsaddr); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

// mkjobfs creates the static portion of the jobd file hierarchy: the 'clone'
// file, and the 'jobs' directory at the root of the hierarchy.
func mkjobfs() (*srv.File, error) {
	var err error

	user := p.OsUsers.Uid2User(os.Geteuid())

	root := new(srv.File)

	err = root.Add(nil, "/", user, nil, p.DMDIR|0555, nil)
	if err != nil {
		return nil, err
	}

	err = mkCloneFile(root, user)
	if err != nil {
		return nil, err
	}

	jobsroot, err = mkJobsDir(root, user)
	if err != nil {
		return nil, err
	}

	return root, nil
}
