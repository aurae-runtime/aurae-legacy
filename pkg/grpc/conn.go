package grpc

import (
	"github.com/libp2p/go-libp2p-core/network"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/sirupsen/logrus"
	"net"
)

var _ net.Conn = &Conn{}

type Conn struct {
	network.Stream
}

// LocalAddr returns the local address.
func (c *Conn) LocalAddr() net.Addr {
	addr := c.Stream.Conn().LocalMultiaddr()
	if na, err := manet.ToNetAddr(addr); err == nil {
		logrus.Infof("[gRPC] Local Addr: %s", na.String())
		return na
	}

	logrus.Warnf("Unable to calculate local address from peer to peer network. %v", addr)
	return &net.TCPAddr{IP: net.ParseIP(""), Port: 0}
}

// RemoteAddr returns the remote address.
func (c *Conn) RemoteAddr() net.Addr {
	addr := c.Stream.Conn().RemoteMultiaddr()
	if na, err := manet.ToNetAddr(addr); err == nil {
		logrus.Infof("[gRPC] Remote Addr: %s", na.String())
		return na
	}

	logrus.Warnf("Unable to calculate remote address from peer to peer network. %v", addr)
	return &net.TCPAddr{IP: net.ParseIP(""), Port: 0}
}
