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
	"syscall"
)

// Root is the root of the auraeFS tree.
// Root represents a single Inode and is the
// base of the filesystem tree for auraeFS.
type Root struct {
	mountpoint string
	fs.Inode
	client *client.Client
}

// NewRoot is where the magic happens.
func NewRoot(mountpoint string, c *client.Client) *Root {
	root := &Root{
		mountpoint: mountpoint,
		client:     c,
	}

	// App

	return root
}

// OnAdd is called when the inode is initialized
func (r *Root) OnAdd(ctx context.Context) {
	// Less is more
}

func (r *Root) NewRegularSubfile(ctx context.Context, name string, data []byte) {
	r.AddChild(name,
		r.NewInode(ctx, NewRegularFile(DefaultauraeFSINodePermissions, data),
			fs.StableAttr{
				Ino:  Ino(),
				Mode: fuse.S_IFREG,
			}), true)
}

func (r *Root) NewRegularSubdirectory(ctx context.Context, name string) {
	r.AddChild(name,
		r.NewInode(ctx, NewRegularDir(name),
			fs.StableAttr{Ino: Ino(),
				Mode: fuse.S_IFDIR}), false)
}

func (r *Root) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = DefaultauraeFSINodePermissions
	return 0
}

func (r *Root) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	var dirents []fuse.DirEntry
	listResp, err := r.client.ListRPC(ctx, &rpc.ListReq{
		Key: "/",
	})
	if err != nil {
		return fs.NewListDirStream(dirents), 1
	}
	for filename, _ := range listResp.Entries {
		dirents = append(dirents, fuse.DirEntry{
			Mode: DefaultauraeFSINodePermissions,
			Name: filename,
			Ino:  Ino(),
		})
	}
	return fs.NewListDirStream(dirents), 0
}

func (r *Root) Opendir(ctx context.Context) syscall.Errno {
	return 0
}

// Root Attributes
var _ fs.NodeGetattrer = &Root{}
var _ fs.NodeOnAdder = &Root{}
var _ fs.NodeOpendirer = &Root{}
var _ fs.NodeReaddirer = &Root{}
