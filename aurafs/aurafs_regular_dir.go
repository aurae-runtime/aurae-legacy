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
	"os"
	"syscall"
)

// Attributes
var _ fs.NodeGetattrer = &RegularDir{}
var _ fs.NodeReaddirer = &RegularDir{}
var _ fs.NodeMkdirer = &RegularDir{}
var _ fs.NodeOpendirer = &RegularDir{}
var _ fs.NodeRmdirer = &RegularDir{}

type RegularDir struct {
	path string
	mode uint32
	fs.Inode
	client *client.Client
}

func NewRegularDir(c *client.Client, path string) *RegularDir {
	s, err := os.Stat(path)
	if err != nil {
		return &RegularDir{
			path:   path,
			mode:   DefaultauraeFSINodePermissions,
			Inode:  fs.Inode{},
			client: c,
		}
	}
	return &RegularDir{
		path:   path,
		mode:   uint32(s.Mode()),
		Inode:  fs.Inode{},
		client: c,
	}
}

func (r *RegularDir) Opendir(ctx context.Context) syscall.Errno {
	return 0
}

func (r *RegularDir) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, 0
}

func (r *RegularDir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	if r.client == nil {
		return fs.NewListDirStream(dirents), 0
	}
	listResp, err := r.client.ListRPC(ctx, &rpc.ListReq{
		Key: r.path,
	})
	logrus.Infof("%v", listResp.Entries)
	logrus.Infof("path: %s", r.path)
	if err != nil {
		return fs.NewListDirStream(dirents), 0
	}
	for filename, content := range listResp.Entries {
		var mode uint32
		var ino uint64
		if content == "" {
			mode = fuse.S_IFDIR
			ino = r.NewRegularSubdirectory(ctx, r.client, filename)
		} else {
			mode = fuse.S_IFREG
			ino = r.NewRegularSubfile(ctx, r.client, filename, []byte(content))
		}
		dirents = append(dirents, fuse.DirEntry{
			Mode: mode,
			Name: filename,
			Ino:  ino,
		})
	}
	return fs.NewListDirStream(dirents), 0
}

func (r *RegularDir) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = r.mode
	return 0
}

func (r *RegularDir) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, 0, 0
}

func (r *RegularDir) Rmdir(ctx context.Context, name string) syscall.Errno {
	return 0
}

func (r *RegularDir) NewRegularSubfile(ctx context.Context, c *client.Client, name string, data []byte) uint64 {
	i := Ino()
	r.AddChild(name,
		r.NewInode(ctx, NewRegularFile(c, DefaultauraeFSINodePermissions, data),
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFREG,
			}), true)
	return i
}

func (r *RegularDir) NewRegularSubdirectory(ctx context.Context, c *client.Client, name string) uint64 {
	i := Ino()
	r.AddChild(name,
		r.NewInode(ctx, NewRegularDir(c, name),
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFDIR,
			}), false)
	return i
}
