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

package auraefs

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

// Open is how all files are opened
//
//               FOPEN_DIRECT_IO
//                     Bypass page cache for this open file.
//
//              FOPEN_KEEP_CACHE
//                     Don't invalidate the data cache on open.
//
//              FOPEN_NONSEEKABLE
//                     The file is not seekable.
func (f *File) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	logrus.Debugf("%s --[f]--> Open()", f.path)
	return nil, fuse.FOPEN_DIRECT_IO, Okay
}

func (f *File) Write(ctx context.Context, fh fs.FileHandle, data []byte, off int64) (uint32, syscall.Errno) {
	logrus.Debugf("%s --[f]--> Write()", f.path)
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
	logrus.Debugf("%s --[f]--> Getattr()", f.path)
	f.mu.Lock()
	defer f.mu.Unlock()
	out.Attr = f.Attr
	out.Attr.Size = uint64(len(f.Data))
	return Okay
}

func (f *File) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	logrus.Debugf("%s --[f]--> Setattr()", f.path)
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
	logrus.Debugf("%s --[f]--> Flush()", f.path)
	setResp, err := c.SetRPC(ctx, &rpc.SetReq{
		Key: f.path,
		Val: string(f.Data),
	})
	if err != nil {
		logrus.Warningf("Unable to SetRPC on Aurae core daemon: %v", err)
		return 1
	}
	if setResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to SetRPC on Aurae core daemon: %v", setResp)
		return 1
	}
	//f.Data = []byte("") // Reset the file content on Flush() if we need it again we pull it from the server.
	return 0
}

func (f *File) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	logrus.Debugf("%s --[f]--> Read().Len(%d)", f.path, len(f.Data))
	f.mu.Lock()
	defer f.mu.Unlock()

	getResp, err := c.GetRPC(ctx, &rpc.GetReq{
		Key: f.path,
	})
	if err != nil {
		logrus.Warningf("Unable to GetRPC on Aurae core daemon: %v", err)
		return fuse.ReadResultData(f.Data), 0
	}
	if getResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to GetRPC on Aurae core daemon: %v", getResp)
		return fuse.ReadResultData(f.Data), 0
	}
	f.Data = []byte(getResp.Val)

	end := int(off) + len(dest)
	if end > len(f.Data) {
		end = len(f.Data)
	}
	return fuse.ReadResultData(f.Data[off:end]), Okay
}

func (f *File) Unlink(ctx context.Context, name string) syscall.Errno {
	logrus.Debugf("%s --[f]--> Unlink(%s)", f.path, name)
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
