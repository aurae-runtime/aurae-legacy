package client

import (
	"github.com/kris-nova/aurae/rpc"
	p2p "github.com/libp2p/go-libp2p"
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

func (c *Client) Connect() error {
	_, err := p2p.New()
	if err != nil {
		return err
	}
	return nil
}

//func (c *Client) Connect() error {
//	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", c.socket), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
//	if err != nil {
//		return err
//	}
//	client := rpc.NewCoreClient(conn)
//	c.CoreClient = client
//	return nil
//}
