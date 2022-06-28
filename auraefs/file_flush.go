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
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"syscall"
)

var _ fs.NodeFlusher = &File{}

func (f *File) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	logrus.Debugf("%s --[f]--> Flush()", f.path)
	if c == nil {
		return 0
	}
	setResp, err := c.Set(ctx, &rpc.SetReq{
		Key: f.path,
		Val: string(f.Data),
	})
	if err != nil {
		logrus.Warningf("Unable to Set on Aurae core daemon: %v", err)
		return 1
	}
	if setResp.Code != core.CoreCode_OKAY {
		logrus.Warningf("Failure to Set on Aurae core daemon: %v", setResp)
		return 1
	}
	//f.Data = []byte("") // Reset the file content on Flush() if we need it again we pull it from the server.
	return 0
}
