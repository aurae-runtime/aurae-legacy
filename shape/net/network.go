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

package net

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae/rpc"
	"github.com/kris-nova/aurae/shape"
	"github.com/kris-nova/aurae/shape/local"
)

var _ shape.AuraeFilesystem = &NetworkFilesystem{}

// NetworkFilesystem is a simple state store for IPC over unix domain sockets on localhost
type NetworkFilesystem struct {
	rpc.UnimplementedAuraeFilesystemServer
	internal *local.LocalFilesystem
}

func NewNetworkFilesystem(internal *local.LocalFilesystem) *NetworkFilesystem {
	return &NetworkFilesystem{
		internal: internal,
	}
}

func (a *NetworkFilesystem) Destroy() error {
	return a.internal.Destroy()
}

func (a *NetworkFilesystem) Initialize(cfg *shape.Config) error {
	return a.internal.Initialize(cfg)
}

func (a *NetworkFilesystem) MountPoint() string {
	return a.internal.MountPoint()
}

func (a *NetworkFilesystem) Socket() string {
	return a.internal.Socket()
}

func (a *NetworkFilesystem) Set(key, value string) error {
	return a.internal.Set(key, value)
}

func (a *NetworkFilesystem) Get(key string) string {
	return a.internal.Get(key)
}

func (a *NetworkFilesystem) SetRPC(ctx context.Context, set *rpc.SetReq) (*rpc.SetResp, error) {
	err := a.Set(set.Key, set.Val)
	if err != nil {
		return &rpc.SetResp{
			Code: -1,
		}, fmt.Errorf("set() filesystem error: %v", err)
	}
	return &rpc.SetResp{
		Code: 1,
	}, nil
}
func (a *NetworkFilesystem) GetRPC(ctx context.Context, get *rpc.GetReq) (*rpc.GetResp, error) {
	val := a.Get(get.Key)
	return &rpc.GetResp{
		Code: 1,
		Val:  val,
	}, nil
}
