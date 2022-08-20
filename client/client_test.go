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

package client

import (
	"context"
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/pkg/core/local"
	p2pgrpc "github.com/kris-nova/aurae/pkg/grpc"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/rpc/rpc"
	"testing"
)

func TestPeer2PeerConnect(t *testing.T) {
	var err error

	// Peer 1
	p1 := peer.Self()
	err = p1.Establish(context.Background(), 0)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	c1 := NewClient()

	// Peer 2
	p2 := peer.NewPeer(name.New("p2"))
	err = p2.Establish(context.Background(), 1)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	// gRPC Server (listen p2)
	proto := p2pgrpc.NewGRPCProtocol(context.Background(), p2.Host())
	if err != nil {
		t.Errorf("unable to listen p2: %v", err)
	}
	server := proto.GetGRPCServer()
	// server.Serve() is called in NewGRPCProtocol\
	rpc.RegisterCoreServer(server,
		core.NewService(local.NewState("/tmp/aurae.test")))

	// gRPC Client (p1 -> p2)
	err = c1.ConnectPeer(p1, p2.Host().ID())
	if err != nil {
		t.Errorf("unable to connect client p1 -> p2: %v", err)
		t.FailNow()
	}

	_, err = c1.Set(context.Background(), &rpc.SetReq{
		Key: "beeps",
		Val: "boops",
	})
	if err != nil {
		t.Errorf("unable to set key/value over peer to peer grpc: %v", err)
	}

}
