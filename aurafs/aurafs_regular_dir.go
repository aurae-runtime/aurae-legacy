/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package aurafs

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"syscall"
)

// Attributes
var _ fs.NodeGetattrer = &RegularDir{}
var _ fs.NodeReaddirer = &RegularDir{}
var _ fs.NodeMkdirer = &RegularDir{}
var _ fs.NodeOpendirer = &RegularDir{}
var _ fs.NodeRmdirer = &RegularDir{}

type RegularDir struct {
	path string
	mode uint32
	fs.Inode
}

// NewRegularDir will build a passthrough directory based on whatever
// is available on the disk.
func NewRegularDir(path string) *RegularDir {
	path = filepath.Join("/aurae", path)
	logrus.Infof("auraeFS RegularDir: %s", path)
	s, err := os.Stat(path)
	if err != nil {
		return &RegularDir{
			path:  path,
			mode:  DefaultauraeFSINodePermissions,
			Inode: fs.Inode{},
		}
	}
	return &RegularDir{
		path:  path,
		mode:  uint32(s.Mode()),
		Inode: fs.Inode{},
	}
}

func (r *RegularDir) Opendir(ctx context.Context) syscall.Errno {
	return 0
}

func (r *RegularDir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, 0
}

func (r *RegularDir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	return fs.NewListDirStream(dirents), 0
}

func (r *RegularDir) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = r.mode
	return 0
}

func (r *RegularDir) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, 0, 0
}

func (r *RegularDir) Rmdir(ctx context.Context, name string) syscall.Errno {
	return 0
}
