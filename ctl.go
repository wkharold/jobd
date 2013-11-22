package main

import (
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p"
	"github.com/wkharold/jobd/deps/code.google.com/p/go9p/p/srv"
)

type ctlfile struct {
	srv.File
	jobsdir []byte
}

func (cf ctlfile) Read(fid *srv.FFid, buf []byte, offset uint64) (int, error) {
	return 0, nil
}

func (cf *ctlfile) Write(fid *srv.FFid, data []byte, offset uint64) (int, error) {
	return 0, nil
}

func (cf *ctlfile) Wstat(fid *srv.FFid, dir *p.Dir) error {
	return nil
}
