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
	"path/filepath"
	"sync"
	"syscall"
)

// Attributes
var _ fs.FileHandle = &File{}
var _ fs.NodeOpener = &File{}
var _ fs.NodeReader = &File{}
var _ fs.NodeWriter = &File{}
var _ fs.NodeSetattrer = &File{}
var _ fs.NodeFlusher = &File{}
var _ fs.NodeGetattrer = &File{}
var _ fs.NodeUnlinker = &File{}

// Model
var _ fs.MemRegularFile

type File struct {
	fs.Inode

	path string // Absolute path

	mu   sync.Mutex
	Data []byte
	Attr fuse.Attr
}

func NewFile(path string, data []byte) *File {
	return &File{
		Inode: fs.Inode{},
		path:  path,
		mu:    sync.Mutex{},
		Data:  data,
		Attr: fuse.Attr{
			Ino:  Ino(),
			Mode: ModeRWOnly,
		}, // Set default attributes here
	}
}

func (f *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, fuse.FOPEN_KEEP_CACHE, Okay
}

func (f *File) Write(ctx context.Context, fh fs.FileHandle, data []byte, off int64) (uint32, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()
	end := int64(len(data)) + off
	if int64(len(f.Data)) < end {
		n := make([]byte, end)
		copy(n, f.Data)
		f.Data = n
	}

	copy(f.Data[off:off+int64(len(data))], data)

	return uint32(len(data)), 0
}

func (f *File) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	f.mu.Lock()
	defer f.mu.Unlock()
	out.Attr = f.Attr
	out.Attr.Size = uint64(len(f.Data))
	return Okay
}

func (f *File) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	f.mu.Lock()
	defer f.mu.Unlock()
	if sz, ok := in.GetSize(); ok {
		f.Data = f.Data[:sz]
	}
	out.Attr = f.Attr
	out.Size = uint64(len(f.Data))
	return Okay
}

func (f *File) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	return 0
}

func (f *File) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	f.mu.Lock()
	defer f.mu.Unlock()
	end := int(off) + len(dest)
	if end > len(f.Data) {
		end = len(f.Data)
	}
	return fuse.ReadResultData(f.Data[off:end]), Okay
}

func (f *File) Unlink(ctx context.Context, name string) syscall.Errno {
	logrus.Debugf("%s -> Unlink(%s)", f.path, name)
	rmResp, err := c.RemoveRPC(ctx, &rpc.RemoveReq{
		Key: filepath.Join(f.path, name),
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
