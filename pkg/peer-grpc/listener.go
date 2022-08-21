package grpc

import (
	"context"
	manet "github.com/multiformats/go-multiaddr/net"
	"github.com/sirupsen/logrus"
	"io"
	"net"
)

// grpcListener implements the net.Listener interface.
type grpcListener struct {
	*GRPCProtocol
	listenerCtx       context.Context
	listenerCtxCancel context.CancelFunc
}

// newGrpcListener builds a new GRPC listener.
func newGrpcListener(proto *GRPCProtocol) net.Listener {
	l := &grpcListener{
		GRPCProtocol: proto,
	}
	l.listenerCtx, l.listenerCtxCancel = context.WithCancel(proto.ctx)
	return l
}

// Accept waits for and returns the next connection to the listener.
func (l *grpcListener) Accept() (net.Conn, error) {
	select {
	case <-l.listenerCtx.Done():
		return nil, io.EOF
	case stream := <-l.streamCh:
		return &Conn{Stream: stream}, nil
	}
}

// Addr returns the listener's network address.
func (l *grpcListener) Addr() net.Addr {
	listenAddrs := l.host.Network().ListenAddresses()
	if len(listenAddrs) > 0 {
		for _, addr := range listenAddrs {
			if na, err := manet.ToNetAddr(addr); err == nil {
				return na
			}
		}
	}
	logrus.Warnf("Unable to calculate listen address from peer to peer network. %v", listenAddrs)
	return &net.TCPAddr{IP: net.ParseIP(""), Port: 0}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *grpcListener) Close() error {
	l.listenerCtxCancel()
	return nil
}
