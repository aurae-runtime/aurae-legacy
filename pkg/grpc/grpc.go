package grpc

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

// Protocol is the GRPC-over-libp2p protocol.
const Protocol protocol.ID = "/grpc/0.0.1"

// GRPCProtocol is the GRPC-transported protocol handler.
type GRPCProtocol struct {
	ctx        context.Context
	host       host.Host
	grpcServer *grpc.Server
	streamCh   chan network.Stream
}

// NewGRPCProtocol attaches the GRPC protocol to a host.
func NewGRPCProtocol(ctx context.Context, host host.Host) *GRPCProtocol {
	grpcServer := grpc.NewServer()
	grpcProtocol := &GRPCProtocol{
		ctx:        ctx,
		host:       host,
		grpcServer: grpcServer,
		streamCh:   make(chan network.Stream),
	}
	host.SetStreamHandler(Protocol, grpcProtocol.HandleStream)
	// Serve will not return until Accept fails, when the ctx is canceled.
	go func() {
		err := grpcServer.Serve(newGrpcListener(grpcProtocol))
		if err != nil {
			logrus.Errorf("unable to start peer grpc server: %v", err)
		}
	}()
	return grpcProtocol
}

// GetGRPCServer returns the grpc server.
func (p *GRPCProtocol) GetGRPCServer() *grpc.Server {
	return p.grpcServer
}

// HandleStream handles an incoming stream.
func (p *GRPCProtocol) HandleStream(stream network.Stream) {
	//logrus.Infof("[grpc] Stream! %s", stream.ID())
	select {
	case <-p.ctx.Done():
		return
	case p.streamCh <- stream:
	}
}
