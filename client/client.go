package client

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/rpc"
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
}

// NewClient will  only be able to authenticate with a local socket.
//
// This mechanism and simple guarantee is what will enable the system to
// operate securely will offline and on the edge.
//
// After a local client has been authenticated the functionality to leverage
// the internal Aurae peering and routing mechanisms are now available.
//
// An authenticated client can use NewPeer() to connect to a proxy in the
// network.
//
// Both Connect() and NewPeer() return unique instances of the same client
// to the user (if successful).
//
// Clients can be chained together to navigate the Aurae mesh.
func NewClient(socket string) *Client {
	return &Client{
		socket:    socket,
		connected: false,
	}
}

// Connect solves identity authorization.
//
// Connect will read authorization certificate material
// on the local filesystem (if it exists) and attempt to
// authenticate with a local unix domain socket.
//
// TODO manage cert material and fix auth
func (c *Client) Connect() error {
	logrus.Warnf("mTLS disabled. running insecure.")
	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", c.socket), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return err
	}
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

func (c *Client) NewPeer(service string) (*Client, error) {
	logrus.Infof("Creating new proxy: %s", service)
	proxyResp, err := c.LocalProxy(context.Background(), &rpc.LocalProxyReq{})
	if err != nil {
		return nil, fmt.Errorf("unable to proxy: %v", err)
	}
	if proxyResp.Code != core.CoreCode_OKAY {
		return nil, fmt.Errorf("unable to create proxy socket: %s", proxyResp.Message)
	}
	return NewClient(proxyResp.Socket), nil
}
