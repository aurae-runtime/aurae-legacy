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
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/rpc"
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

func (s *Service) PeerRequest(ctx context.Context, in *rpc.PeerRequestReq) (*rpc.PeerRequestResp, error) {

	return &rpc.PeerRequestResp{
		Token: "",
	}, nil
}

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
