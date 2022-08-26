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
	"fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/config"
	"github.com/sirupsen/logrus"
	"path"
	"sync"
)

type Dir struct {
	fs.Inode

	path string
	i    uint64
	mu   sync.Mutex
	mode uint32
	Attr fuse.Attr
}

func NewDir(path string) *Dir {
	logrus.Debugf("NewDir: %s", path)
	var i uint64
	i = Ino()
	return &Dir{
		i:     i,
		path:  path,
		Inode: fs.Inode{},
		mu:    sync.Mutex{},
		Attr: fuse.Attr{
			Ino:  i,
			Mode: ModeX,
		}, // Set default attributes here
	}
}

func (n *Dir) NewSubFile(ctx context.Context, name string, data []byte) (uint64, *File) {
	i := Ino()
	if c == nil {
		return 0, nil
	}
	setResp, err := c.Set(ctx, &aurae.SetReq{
		Key: name, // No trailing slash (file)
	})
	if err != nil {
		logrus.Warningf("Unable to Set on Aurae config daemon: %v", err)
		return 0, nil
	}
	if setResp.Code != config.CoreCode_OKAY {
		logrus.Warningf("Failure to Set on Aurae config daemon: %v", setResp)
		return 0, nil
	}

	file := NewFile(path.Join(n.path, name), data)
	n.AddChild(name,
		n.NewInode(ctx, file,
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFREG,
			}), true)
	return i, file
}

func (n *Dir) NewSubDir(ctx context.Context, name string) (uint64, *Dir) {
	i := Ino()
	if c == nil {
		return 0, nil
	}
	setResp, err := c.Set(ctx, &aurae.SetReq{
		Key: fmt.Sprintf("%s/", name), // Trailing slash (dir)
	})
	if err != nil {
		logrus.Warningf("Unable to Set on Aurae config daemon: %v", err)
		return 0, nil
	}
	if setResp.Code != config.CoreCode_OKAY {
		logrus.Warningf("Failure to Set on Aurae config daemon: %v", setResp)
		return 0, nil
	}

	dir := NewDir(path.Join(n.path, name))
	n.AddChild(name,
		n.NewInode(ctx, dir,
			fs.StableAttr{
				Ino:  i,
				Mode: fuse.S_IFDIR,
			}), true)
	return i, dir
}
