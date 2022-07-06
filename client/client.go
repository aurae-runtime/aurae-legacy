package client

import (
	"fmt"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	rpc.CoreClient
	rpc.RuntimeClient
	rpc.ScheduleClient
	//rpc.ProxyClient

	socket    string
	connected bool
	peer      *peer.Peer
}

func NewClient() *Client {
	return &Client{
		connected: false,
	}
}

//func (c *Client) ConnectPeer(self *peer.Peer, to peer2.ID) error {
//	// Cache the peer
//	c.peer = self
//
//	grpcProto := p2pgrpc.NewGRPCProtocol(context.Background(), &self.RHost)
//	conn, err := grpcProto.Dial(context.Background(), to, grpc.WithTimeout(time.Second*6), grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		return err
//	}
//	return c.establish(conn)
//}

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
	//proxy := rpc.NewProxyClient(conn)
	//c.ProxyClient = proxy
	c.connected = true
	logrus.Warnf("Connected to grpc")

	return nil
}

//func (c *Client) NewPeer(service string) (*Client, error) {
//	logrus.Infof("Creating new proxy: %s", service)
//	proxyResp, err := c.LocalProxy(context.Background(), &rpc.LocalProxyReq{})
//	if err != nil {
//		return nil, fmt.Errorf("unable to proxy: %v", err)
//	}
//	if proxyResp.Code != core.CoreCode_OKAY {
//		return nil, fmt.Errorf("unable to create proxy socket: %s", proxyResp.Message)
//	}
//	return NewClient(proxyResp.Socket), nil
//}
