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
