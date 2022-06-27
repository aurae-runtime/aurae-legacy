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
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"sync"
	"syscall"
)

// Attributes
var _ fs.NodeGetattrer = &Dir{}
var _ fs.NodeOnAdder = &Dir{}
var _ fs.NodeOpendirer = &Dir{}
var _ fs.NodeReaddirer = &Dir{}
var _ fs.NodeMkdirer = &Dir{}
var _ fs.NodeRmdirer = &Dir{}

var _ fs.MemRegularFile

type Dir struct {
	fs.Inode

	path string // Absolute path
	i    uint64
	mu   sync.Mutex
	mode uint32
	Attr fuse.Attr
}

func NewDir(path string) *Dir {
	logrus.Debugf("NewDir: %s", path)
	var i uint64
	i = Ino()
	return &Dir{
		i:     i,
		path:  path,
		Inode: fs.Inode{},
		mu:    sync.Mutex{},
		Attr: fuse.Attr{
			Ino:  i,
			Mode: ModeX,
		}, // Set default attributes here
	}
}

func (n *Dir) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	logrus.Debugf("%s -> Getattr()", n.path)
	n.mu.Lock()
	defer n.mu.Unlock()
	out.Mode = n.mode
	out.Ino = n.i
	return 0
}

func (n *Dir) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	n.mu.Lock()
	defer n.mu.Unlock()
	logrus.Debugf("%s -> Setattr()", n.path)
	return Okay
}

func (n *Dir) Opendir(ctx context.Context) syscall.Errno {
	logrus.Debugf("%s -> Opendir()", n.path)
	return 0
}

func (n *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	logrus.Debugf("%s -> Mkdir(%s)", n.path, name)
	_, dir := n.NewSubDir(ctx, path.Join(n.path, name))
	return &dir.Inode, 0
}

func (n *Dir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	logrus.Debugf("%s -> Readdir()", n.path)
	if c == nil {
		return fs.NewListDirStream(dirents), 0
	}
	logrus.Debugf("dir.Readdir() -> client.ListRPC() path=%s", n.path)
	listResp, err := c.ListRPC(ctx, &rpc.ListReq{
		Key: n.path,
	})
	if err != nil {
		return fs.NewListDirStream(dirents), 0
	}
	for filename, node := range listResp.Entries {
		var mode uint32
		var ino uint64
		if node.GetFile() {
			mode = fuse.S_IFREG
			getResp, err := c.GetRPC(ctx, &rpc.GetReq{
				Key: filename,
			})
			if err != nil {
				return fs.NewListDirStream(dirents), 0
			}
			ino, _ = n.NewSubFile(ctx, filename, []byte(getResp.Val))
		} else {
			mode = fuse.S_IFDIR
			ino, _ = n.NewSubDir(ctx, filename)
		}
		dirents = append(dirents, fuse.DirEntry{
			Mode: mode,
			Name: filename,
			Ino:  ino,
		})
	}
	return fs.NewListDirStream(dirents), 0
}

func (n *Dir) Rmdir(ctx context.Context, name string) syscall.Errno {
	logrus.Debugf("%s -> Rmdir()", n.path)
	return 0
}

func (n *Dir) OnAdd(ctx context.Context) {
	logrus.Debugf("%s -> OnAdd()", n.path)
	// Less is more
}

func (n *Dir) NewSubFile(ctx context.Context, name string, data []byte) (uint64, *File) {
	i := Ino()
	file := NewFile(path.Join(n.path, name), data)
	n.AddChild(name,
		n.NewInode(ctx, file,
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFREG,
			}), true)
	return i, file
}

func (n *Dir) NewSubDir(ctx context.Context, name string) (uint64, *Dir) {
	i := Ino()
	setResp, err := c.SetRPC(ctx, &rpc.SetReq{
		Key: filepath.Join(name, "/"),
	})
	if err != nil {
		logrus.Warningf("Unable to SetRPC on Aurae core daemon: %v", err)
		return 0, nil
	}
	if setResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to SetRPC on Aurae core daemon: %v", setResp)
		return 0, nil
	}

	dir := NewDir(path.Join(n.path, name))
	n.AddChild(name,
		n.NewInode(ctx, dir,
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFDIR,
			}), true)
	return i, dir
}
