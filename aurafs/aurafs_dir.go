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
	"fmt"
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
var _ fs.NodeCreater = &Dir{}
var _ fs.NodeUnlinker = &Dir{}

var _ fs.MemRegularFile

type Dir struct {
	fs.Inode

	path string // Absolute path
	i    uint64
	mu   sync.Mutex
	mode uint32
	Attr fuse.Attr
	Data []byte
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
	logrus.Debugf("%s --[d]--> Getattr()", n.path)
	n.mu.Lock()
	defer n.mu.Unlock()
	out.Mode = n.mode
	out.Ino = n.i
	return 0
}

func (n *Dir) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	n.mu.Lock()
	defer n.mu.Unlock()
	logrus.Debugf("%s --[d]--> Setattr()", n.path)
	return Okay
}

func (n *Dir) Opendir(ctx context.Context) syscall.Errno {
	logrus.Debugf("%s --[d]--> Opendir()", n.path)
	return 0
}

func (n *Dir) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (node *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	logrus.Debugf("%s --[d]--> Create(%s)", n.path, name)
	_, file := n.NewSubFile(ctx, path.Join(n.path, name), []byte("")) // touch
	return &file.Inode, nil, 0, 0
}

func (n *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	logrus.Debugf("%s --[d]--> Mkdir(%s)", n.path, name)
	_, dir := n.NewSubDir(ctx, path.Join(n.path, name))
	return &dir.Inode, 0
}

func (n *Dir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	logrus.Debugf("%s --[d]--> Readdir()", n.path)
	if c == nil {
		return fs.NewListDirStream(dirents), 0
	}
	logrus.Debugf("dir.Readdir() --[d]--> client.ListRPC() path=%s", n.path)
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
	logrus.Debugf("%s --[d]--> Rmdir()", n.path)
	rmResp, err := c.RemoveRPC(ctx, &rpc.RemoveReq{
		Key: name,
	})
	if err != nil {
		logrus.Warningf("Unable to RemoveRPC on Aurae core daemon: %v", err)
		return 0
	}
	if rmResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to RemoveRPC on Aurae core daemon: %v", rmResp)
		return 0
	}
	return 0
}

func (n *Dir) OnAdd(ctx context.Context) {
	logrus.Debugf("%s --[d]--> OnAdd()", n.path)
	// Less is more
}

func (n *Dir) NewSubFile(ctx context.Context, name string, data []byte) (uint64, *File) {
	i := Ino()
	setResp, err := c.SetRPC(ctx, &rpc.SetReq{
		Key: name, // No trailing slash (file)
	})
	if err != nil {
		logrus.Warningf("Unable to SetRPC on Aurae core daemon: %v", err)
		return 0, nil
	}
	if setResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to SetRPC on Aurae core daemon: %v", setResp)
		return 0, nil
	}

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
		Key: fmt.Sprintf("%s/", name), // Trailing slash (dir)
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

func (n *Dir) Unlink(ctx context.Context, name string) syscall.Errno {
	logrus.Debugf("%s --[d]--> Unlink(%s)", n.path, name)
	rmResp, err := c.RemoveRPC(ctx, &rpc.RemoveReq{
		Key: filepath.Join(n.path, name),
	})
	if err != nil {
		logrus.Warningf("Unable to RemoveRPC on Aurae core daemon: %v", err)
		return 0
	}
	if rmResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to RemoveRPC on Aurae core daemon: %v", rmResp)
		return 0
	}
	return 0
}
