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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/config"
	"github.com/sirupsen/logrus"
	"syscall"
)

var _ fs.NodeFlusher = &File{}

func (f *File) Flush(ctx context.Context, fh fs.FileHandle) syscall.Errno {
	logrus.Debugf("%s --[f]--> Flush()", f.path)
	if c == nil {
		return 0
	}
	setResp, err := c.Set(ctx, &aurae.SetReq{
		Key: f.path,
		Val: string(f.Data),
	})
	if err != nil {
		logrus.Warningf("Unable to Set on Aurae config daemon: %v", err)
		return 1
	}
	if setResp.Code != config.CoreCode_OKAY {
		logrus.Warningf("Failure to Set on Aurae config daemon: %v", setResp)
		return 1
	}
	//f.Data = []byte("") // Reset the file content on Flush() if we need it again we pull it from the server.
	return 0
}
