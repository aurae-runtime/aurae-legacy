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

// mountOptions are fusermount mount options.
//
// These are a combination of mount(8) and fusermount afuse(1)
// mount options.
//
// Most of the generic mount options described in mount are
// supported (ro, rw, suid, nosuid, dev, nodev, exec, noexec, atime,
// noatime, sync, async, dirsync). Filesystems are mounted with
// nodev,nosuid by default, which can only be overridden by a
// privileged user.
//
// More: https://linux.die.net/man/8/mount
// More: https://www.unix.com/man-page/linux/1/fusermount/
var mountOptions = []string{

	// -o nonempty 	      allow mounts over non-empty file/dir
	"nonempty",

	// -o rw 			  mount the filesystem as read/write
	"rw",
}

// Mount will mount the FS and return.
func (a *AuraeFS) Mount() error {

	logrus.Infof("Mountpoint: %s", a.hostMountPath)

	for _, mountOption := range mountOptions {
		logrus.Debugf("Mount option: %s", mountOption)
	}

	// Mount root and its children
	svc, err := fs.Mount(a.hostMountPath, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			Options: mountOptions,
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
			logrus.Warningf("Unable to unmount: %v", err)
			time.Sleep(time.Second * 2)
			err = a.service.Unmount()
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
