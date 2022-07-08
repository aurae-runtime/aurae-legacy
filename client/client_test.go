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
	p2pgrpc "github.com/kris-nova/aurae/pkg/grpc"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/runtime"
	"github.com/kris-nova/aurae/rpc"
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
	// server.Serve() is called in NewGRPCProtocol
	rpc.RegisterRuntimeServer(server, runtime.NewService())

	// Failing here
	// Idea 1: Check if we need to server.Serve() the gRPC server
	// Idea 2: Check if the port number offset needs to be plumbed through
	// Idea 3: Check the net.Conn and net.Dialer used in handshake()

	// gRPC Client (p1 -> p2)
	err = c1.ConnectPeer(p1, p2.Host().ID())
	if err != nil {
		t.Errorf("unable to connect client p1 -> p2: %v", err)
	}

	_, err = c1.Set(context.Background(), &rpc.SetReq{
		Key: "beeps",
		Val: "boops",
	})
	if err != nil {
		t.Errorf("unable to set key/value over peer to peer grpc: %v", err)
	}

}