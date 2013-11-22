package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"

	"flag"
	"os"
)

func main() {
	flfsaddr := flag.String("fsaddr", "0.0.0.0:5640", "Address where job file service listens for connections")

	root, err := mkjobfs()
	if err != nil {
		os.Exit(1)
	}

	s := srv.NewFileSrv(root)
	s.Dotu = true
	s.Start(s)

	if err := s.StartNetListener("tcp", *flfsaddr); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func mkjobfs() (*srv.File, error) {
	user := p.OsUsers.Uid2User(os.Geteuid())

	root := new(srv.File)
	if err := root.Add(nil, "/", user, nil, p.DMDIR|0555, nil); err != nil {
		return nil, err
	}

	return root, nil
}
