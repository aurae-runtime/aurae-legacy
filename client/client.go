package client

import (
	"fmt"
	"github.com/kris-nova/aurae/rpc"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	rpc.CoreServiceClient
	socket string
}

func NewClient(socket string) *Client {
	return &Client{
		socket: socket,
	}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", c.socket), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return err
	}
	client := rpc.NewCoreServiceClient(conn)
	c.CoreServiceClient = client
	return nil
}
