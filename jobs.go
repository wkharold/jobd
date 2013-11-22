package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
)

type jobsdir struct {
	srv.File
}

func (jd jobsdir) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	return 0, nil
}

func (jd *jobsdir) Wstat(fid *srv.FFid, dir *p.Dir) error {
	return nil
}
