package p2pgrpc

import (
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/sirupsen/logrus"
	"net"
)

// streamConn represents a net.Conn wrapped to be compatible with net.conn
type streamConn struct {
	network.Stream
}

// LocalAddr returns the local address.
func (c *streamConn) LocalAddr() net.Addr {
	//addr := c.Stream.Conn().LocalMultiaddr()
	ret := &net.IPAddr{}
	// TODO THIS IS BROKEN
	logrus.Warnf("UNSUPPORTED TRANSLATION")

	return ret
}

// RemoteAddr returns the remote address.
func (c *streamConn) RemoteAddr() net.Addr {
	//addr := c.Stream.Conn().LocalMultiaddr()
	ret := &net.IPAddr{}
	// TODO THIS IS BROKEN
	logrus.Warnf("UNSUPPORTED TRANSLATION")

	return ret
}

var _ net.Conn = &streamConn{}
