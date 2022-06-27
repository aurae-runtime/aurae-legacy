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
	"github.com/kris-nova/aurae/pkg/core/memfs"
	"github.com/kris-nova/aurae/rpc"
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
	Node *memfs.Node

	path string // Absolute path

	mu   sync.Mutex
	mode uint32
	Attr fuse.Attr
}

func NewDir(path string) *Dir {
	var node *memfs.Node
	node = root.Node.AddSubNode(path, "") // TODO We need mkdir!
	return &Dir{
		path:  path,
		Inode: fs.Inode{},
		mu:    sync.Mutex{},
		Attr:  fuse.Attr{}, // Set default attributes here
		Node:  node,
	}
}

func (n *Dir) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, fuse.FOPEN_KEEP_CACHE, Okay
}

func (n *Dir) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	n.mu.Lock()
	defer n.mu.Unlock()
	out.Attr = n.Attr
	return Okay
}

func (n *Dir) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	n.mu.Lock()
	defer n.mu.Unlock()
	out.Attr = n.Attr
	return Okay
}

func (n *Dir) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	return 0
}

func (n *Dir) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n.mu.Lock()
	defer n.mu.Unlock()
	return fuse.ReadResultData([]byte("")), Okay
}

func (n *Dir) Opendir(ctx context.Context) syscall.Errno {
	return 0
}

func (n *Dir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, 0
}

func (n *Dir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	if c == nil {
		return fs.NewListDirStream(dirents), 0
	}
	listResp, err := c.ListRPC(ctx, &rpc.ListReq{
		Key: n.path,
	})
	if err != nil {
		return fs.NewListDirStream(dirents), 0
	}
	for filename, node := range listResp.Entries {
		var mode uint32
		if node.GetFile() {
			mode = fuse.S_IFREG
		} else {
			mode = fuse.S_IFDIR
		}
		dirents = append(dirents, fuse.DirEntry{
			Mode: mode,
			Name: filename,
		})
	}
	return fs.NewListDirStream(dirents), 0
}

func (n *Dir) Rmdir(ctx context.Context, name string) syscall.Errno {
	return 0
}

func (n *Dir) OnAdd(ctx context.Context) {
	// Less is more
}
