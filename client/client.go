package client

import (
	"fmt"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	rpc.CoreClient
	socket string
}

// NewClient will  only be able to authenticate with a local socket.
//
// This mechanism and simple guarantee is what will enable the system to
// operate securely will offline and on the edge.
//
// After a local client has been authenticated the functionality to leverage
// the internal Aurae peering and routing mechanisms are now available.
//
// An authenticated client can use PeerConnect() to connect to a peer in the
// network.
//
// Both Connect() and PeerConnect() return unique instances of the same client
// to the user (if successful).
//
// Clients can be chained together to navigate the Aurae mesh.
func NewClient(socket string) *Client {
	return &Client{
		socket: socket,
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
	client := rpc.NewCoreClient(conn)
	c.CoreClient = client
	return nil
}

func (c *Client) PeerConnect(hostname string) (*Client, error) {
	logrus.Infof("Peer connection: %s", hostname)
	logrus.Warnf("peer connection unsupported. returning local aurae client")
	return c, nil
}
