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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/peer"
	p2pgrpc "github.com/kris-nova/aurae/pkg/peer-grpc"
	peer2peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"time"
)

type Client struct {
	aurae.ConfigClient
	aurae.LocalRuntimeClient
	aurae.SandboxRuntimeClient
	aurae.RegisterClient
	aurae.SystemClient
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
		return fmt.Errorf("unable to initialize required handshake before peer-grpc: %v", err)
	}
	logrus.Debugf("Connecting (gRPC) to: %s...", to.String())
	ctx := context.Background()

	//availableProtocols, err := self.Host().Peerstore().GetProtocols(to)
	//if err != nil {
	//	return fmt.Errorf("unable to list protocols of remote: %v", err)
	//}
	//if len(availableProtocols) < 1 {
	//	return fmt.Errorf("no remote protocols found on remote: %s", to.String())
	//}
	//logrus.Infof("Known remote protocols:")
	//for _, ap := range availableProtocols {
	//	logrus.Infof(" - [%s]", ap)
	//}

	//knownPeerAddrs := self.Host().Peerstore().Addrs(to)
	//if len(knownPeerAddrs) < 1 {
	//	return fmt.Errorf("no multi addrs found on remote: %s", to.String())
	//}
	//logrus.Infof("Known peer multi addresses:")
	//for _, ma := range knownPeerAddrs {
	//	logrus.Infof(" - [%s]", ma.String())
	//}

	grpcProto := p2pgrpc.NewGRPCProtocol(ctx, self.Host())
	conn, err := grpcProto.Dial(ctx, to, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3), grpc.WithBlock())
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

	stat, err := os.Stat(c.socket)
	if err != nil {
		return fmt.Errorf("unable to stat socket: %v", err)
	}
	if stat.Mode()&os.ModeSocket == 0 {
		return fmt.Errorf("file is not of type unix socket")
	}

	logrus.Warnf("mTLS disabled. running insecure.")
	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", c.socket),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return err
	}
	if conn == nil {
		return fmt.Errorf("unable to establish connection")
	}
	return c.establish(conn)
}

func (c *Client) establish(conn grpc.ClientConnInterface) error {
	// Establish the connection from the conn
	cfg := aurae.NewConfigClient(conn)
	c.ConfigClient = cfg
	local := aurae.NewLocalRuntimeClient(conn)
	c.LocalRuntimeClient = local
	sandbox := aurae.NewSandboxRuntimeClient(conn)
	c.SandboxRuntimeClient = sandbox
	register := aurae.NewRegisterClient(conn)
	c.RegisterClient = register
	system := aurae.NewSystemClient(conn)
	c.SystemClient = system
	c.connected = true
	return nil
}
