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
	"time"
)

type AuraeFS struct {
	runtime       bool
	service       *fuse.Server
	hostMountPath string
}

var single bool = false

func NewAuraeFS(mountpoint string, establishedClient *client.Client) *AuraeFS {
	c = establishedClient
	if single {
		panic("AuraeFS already initialized. Unable to initialize multiple AuraeFS within the same thread pool!")
	}
	single = true
	return &AuraeFS{
		runtime:       true,
		hostMountPath: mountpoint,
	}
}

// Mount will mount the FS and return.
func (a *AuraeFS) Mount() error {

	logrus.Infof("Mountpoint: %s", a.hostMountPath)

	// Mount root and its children
	svc, err := fs.Mount(a.hostMountPath, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			Options: []string{"nonempty"},
		},
	})
	if err != nil {
		return err
	}
	a.service = svc
	go func() {
		logrus.Infof("Starting FUSE loop...")
		a.service.Wait() // Run this in a Goroutine so we can respond later
	}()
	logrus.Infof("AuraeFS Started!")
	return nil
}

func (a *AuraeFS) Runtime() {
	quitCh := posix.SignalHandler()
	go func() {
		<-quitCh
		err := a.service.Unmount()
		for err != nil {
			err = a.service.Unmount()
			logrus.Warningf("Unable to unmount: %v", err)
			time.Sleep(time.Second * 2)
		}
		a.runtime = false
	}()

	for a.runtime {
		// TODO Runtime logic
	}
}

var ino uint64 = 1 // 1 is always reserved, and we immediately begin indexing, thus we actually start with 2

func Ino() uint64 {
	ino = ino + 1
	return ino
}
