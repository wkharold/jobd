package main

import (
	"bufio"
	"flag"
	"os"
	"path"
	"strings"

	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
	"github.com/wkharold/jobd/deps/github.com/golang/glog"
)

// jobsroot is the root of the jobd file hierarchy
var jobsroot *jobsdir

// jobsdb is the path to the jobs database
var jobsdb string

func main() {
	flfsaddr := flag.String("fsaddr", "0.0.0.0:5640", "Address where job file service listens for connections")
	fldbdir := flag.String("dbdir", "/var/lib/jobd", "Location of the jobd jobs database")
	fldebug := flag.Bool("debug", false, "9p debugging to stderr")
	flag.Parse()

	var err error

	jobsdb, err = mkjobdb(*fldbdir)
	if err != nil {
		os.Exit(1)
	}

	root, err := mkjobfs()
	if err != nil {
		os.Exit(1)
	}

	switch db, err := os.Open(jobsdb); {
	case err != nil:
		os.Exit(1)
	default:
		scanner := bufio.NewScanner(db)
		for scanner.Scan() {
			data := scanner.Text()
			jdparts := strings.Split(data, ":")
			if len(jdparts) != 3 {
				glog.Errorf("jobdb corruption: invalid job definition (%v)", data)
				os.Exit(1)
			}

			jd, err := mkJobDefinition(jdparts[0], jdparts[1], jdparts[2])
			if err != nil {
				glog.Errorf("unable to create job definition (%v)", err)
				os.Exit(1)
			}

			if err := jobsroot.addJob(*jd); err != nil {
				glog.Errorf("can't add job (%v)", err)
				os.Exit(1)
			}
		}
	}

	s := srv.NewFileSrv(root)
	s.Dotu = true
	if *fldebug {
		s.Debuglevel = 1
	}
	s.Start(s)

	if err := s.StartNetListener("tcp", *flfsaddr); err != nil {
		glog.Errorf("listener failed to start (%v)", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// mkjobdb checks to see if the specified path to the jobd database exists and creates it
// if necessary, it also creates an empty database if none exists and returns the full
// path to the jobs database
func mkjobdb(dbdir string) (string, error) {
	if err := os.MkdirAll(dbdir, 0755); err != nil {
		return "", err
	}

	dbpath := path.Join(dbdir, "jobs.db")

	f, err := os.OpenFile(dbpath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return "", err
	}
	f.Close()

	return dbpath, nil
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
