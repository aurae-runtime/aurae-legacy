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
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/posix"
	"github.com/sirupsen/logrus"
)

// Ensure that Root is of type ModularFilesystem, so we can rip this out later
var _ ModularFilesystem = &Root{}

const (
	RuntimeUnmountTimeoutSeconds   int    = 35
	DefaultauraeFSINodePermissions uint32 = 0755
)

type AuraeFS struct {
	root    *Root
	runtime bool
	service *fuse.Server
	client  *client.Client
}

func NewAuraeFS(mountpoint string, c *client.Client) *AuraeFS {
	return &AuraeFS{
		runtime: true,
		root:    NewRoot(mountpoint, c),
		client:  c,
	}
}

// Mount will mount the FS and return.
func (a *AuraeFS) Mount() error {

	logrus.Infof("Mountpoint: %s", a.root.mountpoint)

	// Mount root and its children
	svc, err := fs.Mount(a.root.mountpoint, a.root, &fs.Options{
		MountOptions: fuse.MountOptions{
			Options: []string{"nonempty"},
		},
	})
	if err != nil {
		return err
	}
	a.service = svc
	logrus.Infof("AuraeFS Started!")
	return nil
}

func (a *AuraeFS) Runtime() {
	quitCh := posix.SignalHandler()
	go func() {
		a.runtime = <-quitCh
	}()
	a.service.Wait()
}

var ino uint64 = 1 // 1 is always reserved, and we immediately begin indexing, thus we actually start with 2

func Ino() uint64 {
	ino = ino + 1
	return ino
}
