/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package proxy

import (
	"context"
	crand "crypto/rand"
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/rpc"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"io"
	mrand "math/rand"
	"time"
)

var _ rpc.ProxyServer = &Service{}

// Service is the main proxy name for managing peer-to-peer connections.
//
// This will implement the proxy server methods defined in rpc/proxy.proto
type Service struct {
	newRandomReader newRandomReader
	rpc.UnimplementedProxyServer
}

// PeerRequest will be a request from a local client for this aurae daemon to begin accepting
// peer requests.
func (s *Service) PeerRequest(ctx context.Context, in *rpc.PeerRequestReq) (*rpc.PeerRequestResp, error) {

	// First create a new key pair from the initialized random seed
	r := s.newRandomReader()
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, fmt.Errorf("unable to generate key pair: %v", err)
	}

	// Default peering options
	//
	// TODO we will most likely need to pull these out into a defaults package

	listenPort := -1 // TODO This will go away once we have domain sockets working

	opts := []p2p.Option{
		p2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		p2p.Identity(priv),
		p2p.DisableRelay(),
	}

	// p2p begin peering.
	host, err := p2p.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to begin peering host: %v", err)
	}
	hostAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID().Pretty()))
	if err != nil {
		return nil, fmt.Errorf("unable to get host addr: %v", err)
	}
	addr := host.Addrs()[0]
	token := addr.Encapsulate(hostAddr).String()
	ids := protocol.ConvertFromStrings([]string{fmt.Sprintf("/aurae/%s", aurae.Version)})
	protocolID := ids[0]

	// Set the stream handler
	logrus.Infof("Setting handler: %s", "HandlePeerConnectStream")
	host.SetStreamHandler(protocolID, peer.HandlePeerConnectStream)

	return &rpc.PeerRequestResp{
		Token: token,
	}, nil
}

// LocalProxy will configure a local unix domain socket that serves as a client proxy.
//
// Once LocalProxy is successful, a new unix domain socket will be available on the filesystem.
// A client can connect to this socket and begin interacting with a remote client using the same
// client interface available on this instance.
//
// Think of LocalProxy as small unix domain socket doors to other nodes.
//
// There's some doors in this house.
func (s *Service) LocalProxy(ctx context.Context, in *rpc.LocalProxyReq) (*rpc.LocalProxyResp, error) {

	logrus.Infof("Peering with name: %s", in.Hostname)
	logrus.Infof("Peering with token: %s", in.Token)

	ret := &rpc.LocalProxyResp{
		Socket:   "/run/todo",
		Hostname: in.Hostname,
		Code:     core.CoreCode_REJECT,
		Message:  "LocalProxy TODO",
	}
	return ret, nil
}

// NewService will need to seed our cryptography libraries.
func NewService() *Service {
	return &Service{
		//reader: randomReaderSeed,
		//reader: randomReaderTimeNow,
		newRandomReader: randomReaderCrypto,
	}
}

type newRandomReader func() io.Reader

func randomReaderCrypto() io.Reader {
	return crand.Reader
}

func randomReaderTimeNow() io.Reader {
	return mrand.New(mrand.NewSource(time.Now().Unix()))
}

var seed int64 = time.Now().Unix()

func randomReaderSeed() io.Reader {
	return mrand.New(mrand.NewSource(seed))
}
