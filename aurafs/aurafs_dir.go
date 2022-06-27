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
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"path"
	"sync"
	"syscall"
)

// Attributes
var _ fs.NodeGetattrer = &Dir{}
var _ fs.NodeOnAdder = &Dir{}
var _ fs.NodeOpendirer = &Dir{}
var _ fs.NodeReaddirer = &Dir{}

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

func (n *Dir) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	logrus.Debugf("%s -> Open()", n.path)
	return nil, fuse.FOPEN_KEEP_CACHE, Okay
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

func (n *Dir) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	logrus.Debugf("%s -> Flush()", n.path)
	return 0
}

func (n *Dir) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n.mu.Lock()
	defer n.mu.Unlock()
	logrus.Debugf("%s -> Read()", n.path)
	return fuse.ReadResultData([]byte("")), Okay
}

func (n *Dir) Opendir(ctx context.Context) syscall.Errno {
	logrus.Debugf("%s -> Opendir()", n.path)
	return 0
}

func (n *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	logrus.Debugf("%s -> Mkdir()", n.path)
	return nil, 0
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
			ino = n.NewSubFile(ctx, filename, []byte(getResp.Val))
		} else {
			mode = fuse.S_IFDIR
			ino = n.NewSubDir(ctx, filename)
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

func (n *Dir) NewSubFile(ctx context.Context, name string, data []byte) uint64 {
	i := Ino()
	n.AddChild(name,
		n.NewInode(ctx, NewFile(path.Join(n.path, name), data),
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFREG,
			}), true)
	return i
}

func (n *Dir) NewSubDir(ctx context.Context, name string) uint64 {
	i := Ino()
	n.AddChild(name,
		n.NewInode(ctx, NewDir(path.Join(n.path, name)),
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFDIR,
			}), true)
	return i
}
