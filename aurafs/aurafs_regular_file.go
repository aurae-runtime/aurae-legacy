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
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"syscall"
)

// Attributes
var _ fs.NodeGetattrer = &RegularFile{}
var _ fs.NodeOpener = &RegularFile{}
var _ fs.FileReader = &RegularFile{}
var _ fs.FileHandle = &RegularFile{}
var _ fs.FileFlusher = &RegularFile{}
var _ fs.MemRegularFile

type RegularFile struct {
	fs.LoopbackNode
	mode uint32
	data []byte
	fs.Inode
	path   string
	client *client.Client
}

func (r *RegularFile) Flush(ctx context.Context) syscall.Errno {
	return 0
}

func (r *RegularFile) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	logrus.Infof("Read()")
	resp, err := r.client.GetRPC(ctx, &rpc.GetReq{
		Key: r.path,
	})
	logrus.Infof("path: %s", r.path)
	logrus.Infof("raw: %s", resp.Val)
	if err != nil {
		return fuse.ReadResultData(r.data), 0
	}
	return fuse.ReadResultData([]byte(resp.Val)), 0
}

func (r *RegularFile) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, fuse.FOPEN_KEEP_CACHE, 0
}

func (r *RegularFile) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = r.mode
	//out.Attr.Size = uint64(len(f.Data))
	return 0
}

func NewRegularFile(c *client.Client, mode uint32, path string, data []byte) *RegularFile {
	return &RegularFile{
		mode:   mode,
		client: c,
		Inode:  fs.Inode{},
		data:   data,
		path:   path,
	}
}

// ---
