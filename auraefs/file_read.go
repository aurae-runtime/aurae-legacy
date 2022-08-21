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
	"github.com/kris-nova/aurae/pkg/config"
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/sirupsen/logrus"
	"syscall"
)

var _ fs.NodeReader = &File{}

func (f *File) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	logrus.Debugf("%s --[f]--> Read().Len(%d)", f.path, len(f.Data))
	f.mu.Lock()
	defer f.mu.Unlock()
	if c == nil {
		return fuse.ReadResultData(f.Data), 0
	}
	getResp, err := c.Get(ctx, &rpc.GetReq{
		Key: f.path,
	})
	if err != nil {
		logrus.Warningf("Unable to Get on Aurae config daemon: %v", err)
		return fuse.ReadResultData(f.Data), 0
	}
	if getResp.Code != config.CoreCode_OKAY {
		logrus.Warningf("Failure to Get on Aurae config daemon: %v", getResp)
		return fuse.ReadResultData(f.Data), 0
	}
	f.Data = []byte(getResp.Val)

	end := int(off) + len(dest)
	if end > len(f.Data) {
		end = len(f.Data)
	}
	return fuse.ReadResultData(f.Data[off:end]), Okay
}
