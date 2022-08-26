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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/sirupsen/logrus"
	"syscall"
)

var _ fs.NodeReaddirer = &Dir{}

func (n *Dir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	logrus.Debugf("%s --[d]--> Readdir()", n.path)
	if c == nil {
		return fs.NewListDirStream(dirents), 0
	}
	logrus.Debugf("dir.Readdir() --[d]--> client.List() path=%s", n.path)
	listResp, err := c.List(ctx, &aurae.ListReq{
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
			getResp, err := c.Get(ctx, &aurae.GetReq{
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
