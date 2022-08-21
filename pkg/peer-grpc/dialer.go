package peer_grpc

import (
	"context"
	"net"
	"time"

	peer "github.com/libp2p/go-libp2p-core/peer"
	"google.golang.org/grpc"
)

// GetDialOption returns the WithDialer option to dial via libp2p.
// note: ctx should be the root context.
func (p *GRPCProtocol) GetDialOption(ctx context.Context) grpc.DialOption {
	return grpc.WithDialer(func(peerIdStr string, timeout time.Duration) (net.Conn, error) {
		subCtx, subCtxCancel := context.WithTimeout(ctx, timeout)
		defer subCtxCancel()
		id, err := peer.Decode(peerIdStr)
		if err != nil {
			return nil, err
		}
		err = p.host.Connect(subCtx, peer.AddrInfo{
			ID: id,
		})
		if err != nil {
			return nil, err
		}
		stream, err := p.host.NewStream(ctx, id, Protocol)
		if err != nil {
			return nil, err
		}
		return &Conn{Stream: stream}, nil
	})
}

func (p *GRPCProtocol) Dial(ctx context.Context, peerID peer.ID, dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {

	// Append with Dialer option
	dialOpsPrepended := append([]grpc.DialOption{p.GetDialOption(ctx)}, dialOpts...)
	return grpc.Dial(peerID.Pretty(), dialOpsPrepended...)
	//return peer-grpc.DialContext(ctx, peerID.Pretty(), dialOpsPrepended...)
}
