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
	"fmt"
	p2pgrpc "github.com/kris-nova/aurae/pkg/grpc"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/rpc"
	peer2peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	rpc.CoreClient
	rpc.RuntimeClient
	rpc.ScheduleClient
	rpc.ProxyClient

	socket    string
	connected bool
	peer      *peer.Peer
}

func NewClient() *Client {
	return &Client{
		connected: false,
	}
}

func (c *Client) ConnectPeer(self *peer.Peer, to peer2peer.ID) error {
	err := self.Handshake(to) // Not necessarily *required* but it's a good check for basic connectivity
	if err != nil {
		return fmt.Errorf("unable to initialize required handshake before grpc: %v", err)
	}
	logrus.Infof("Connecting (gRPC) to: %s...", to.String())
	ctx := context.Background()
	grpcProto := p2pgrpc.NewGRPCProtocol(ctx, self.Host())
	ipfsAddr, err := multiaddr.NewMultiaddr(to.String())
	if err != nil {
		return fmt.Errorf("unable to find IPFS multi address to dial: %v", err)
	}
	protocolID, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		return fmt.Errorf("unable to calculate value for protocol: %v", err)
	}
	peerIDToDial, err := peer2peer.Decode(protocolID)
	if err != nil {
		return fmt.Errorf("unable to decode peerIDToDial: %v", err)
	}
	targetPeerAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", peer2peer.Encode(peerIDToDial)))
	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	self.Host().Peerstore().AddAddr(peerIDToDial, targetAddr, peerstore.PermanentAddrTTL)
	logrus.Infof("NewGRPC with host initialized. Dialing %s...", peerIDToDial)
	conn, err := grpcProto.Dial(ctx, peerIDToDial, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("unable to dial: %v", err)
	}
	err = c.establish(conn)
	if err != nil {
		return fmt.Errorf("unable to establish connection: %v", err)
	}
	return nil
}

func (c *Client) ConnectSocket(sock string) error {

	// Cache the socket
	c.socket = sock

	logrus.Warnf("mTLS disabled. running insecure.")
	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", c.socket),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return err
	}
	return c.establish(conn)
}

func (c *Client) establish(conn grpc.ClientConnInterface) error {
	// Establish the connection from the conn
	core := rpc.NewCoreClient(conn)
	c.CoreClient = core
	runtime := rpc.NewRuntimeClient(conn)
	c.RuntimeClient = runtime
	schedule := rpc.NewScheduleClient(conn)
	c.ScheduleClient = schedule
	proxy := rpc.NewProxyClient(conn)
	c.ProxyClient = proxy
	c.connected = true
	return nil
}
