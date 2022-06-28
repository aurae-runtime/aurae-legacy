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
	"path/filepath"
	"syscall"
)

var _ fs.NodeFlusher = &File{}

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
