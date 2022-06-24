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
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
)

// Ensure that Root is of type ModularFilesystem, so we can rip this out later
var _ ModularFilesystem = &Root{}

const (
	RuntimeUnmountTimeoutSeconds   int    = 35
	DefaultauraeFSINodePermissions uint32 = 0755
)

type AuraeFS struct {
	rpc.AuraeFSServer
	root    *Root
	runtime bool
	service *fuse.Server
}

func NewAuraeFS(mountpoint string) *AuraeFS {
	return &AuraeFS{
		runtime: true,
		root:    NewRoot(mountpoint),
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
	go a.service.Wait()
	logrus.Infof("AuraeFS Started!")
	return nil
}

var ino uint64 = 1 // 1 is always reserved, and we immediately begin indexing, thus we actually start with 2

func Ino() uint64 {
	ino = ino + 1
	return ino
}

// -- gRPC Implementation --

func (a *AuraeFS) SetRPC(ctx context.Context, req *rpc.SetReq) (*rpc.SetResp, error) {
	logrus.Infof("Set: %s %s", req.Key, req.Val)
	a.root.NewRegularSubfile(ctx, req.Key, []byte(req.Val)) // Map key, value directly to filename and []byte data
	resp := &rpc.SetResp{}
	return resp, nil
}

func (a *AuraeFS) GetRPC(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	logrus.Infof("Get: %s", req.Key)
	inode := a.root.GetChild(req.Key)
	if inode == nil {
		return &rpc.GetResp{
			Val:  "",
			Code: -1,
		}, nil
	}
	// LEFT OFF HERE
	//rf := inode.(*RegularFile)
	resp := &rpc.GetResp{
		Code: 1,
		Val:  "",
	}
	return resp, nil
}
