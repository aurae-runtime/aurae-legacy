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
	"syscall"
)

// Root is the root of the auraeFS tree.
// Root represents a single Inode and is the
// base of the filesystem tree for auraeFS.
type Root struct {
	mountpoint string
	fs.Inode
}

// NewRoot is where the magic happens.
func NewRoot(mountpoint string) *Root {
	root := &Root{
		mountpoint: mountpoint,
	}

	// App

	return root
}

// OnAdd is called when the inode is initialized
func (r *Root) OnAdd(ctx context.Context) {
	// Less is more
}

func (r *Root) NewRegularSubfile(ctx context.Context, name string) {
	r.AddChild(name,
		r.NewInode(ctx, NewRegularFile(DefaultauraeFSINodePermissions),
			fs.StableAttr{
				Ino:  Ino(),
				Mode: fuse.S_IFREG,
			}), false)
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

// Root Attributes
var _ fs.NodeGetattrer = &Root{}
var _ fs.NodeOnAdder = &Root{}
