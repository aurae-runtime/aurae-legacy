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

	mu   sync.Mutex
	mode uint32
	Attr fuse.Attr
}

func NewDir(attr fuse.Attr) *File {
	return &File{
		Inode: fs.Inode{},
		mu:    sync.Mutex{},
		Attr:  attr,
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
	//if r.client == nil {
	//	return fs.NewListDirStream(dirents), 0
	//}
	//listResp, err := r.client.ListRPC(ctx, &rpc.ListReq{
	//	Key: r.path,
	//})
	//if err != nil {
	//	return fs.NewListDirStream(dirents), 0
	//}
	//for filename, content := range listResp.Entries {
	//	var mode uint32
	//	var ino uint64
	//	if content == "" {
	//		mode = fuse.S_IFDIR
	//		ino = r.NewRegularSubdirectory(ctx, r.client, filename)
	//	} else {
	//		mode = fuse.S_IFREG
	//		ino = r.NewRegularSubfile(ctx, r.client, filename, []byte(content))
	//	}
	//	dirents = append(dirents, fuse.DirEntry{
	//		Mode: mode,
	//		Name: filename,
	//		Ino:  ino,
	//	})
	//}
	return fs.NewListDirStream(dirents), 0
}

func (n *Dir) Rmdir(ctx context.Context, name string) syscall.Errno {
	return 0
}

func (n *Dir) OnAdd(ctx context.Context) {
	// Less is more
}
