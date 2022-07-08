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

package peer

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"
	"github.com/kris-nova/aurae/pkg/common"
	p2pgrpc "github.com/kris-nova/aurae/pkg/grpc"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/rpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

const (
	DefaultGenerateKeyPairBits int = 2048
	DefaultListenPort          int = 8709
	DefaultPeerPort            int = 8708
)

var emptyKey crypto.PrivKey = &crypto.Ed25519PrivateKey{}

type Peer struct {
	uniqKey     crypto.PrivKey
	established bool
	Name        name.Name
	//Peers       map[string]*Peer
	host host.Host
	//DNS         *NameService
	//internalDNS mdns.Service

	localSocket string
	rpc.CoreClient
	rpc.RuntimeClient
	rpc.ScheduleClient
	rpc.ProxyClient
}

func NewPeer(n name.Name) *Peer {
	golog.SetupLogging(golog.Config{
		Stdout: true,
		Stderr: false,
	})
	golog.SetAllLoggers(golog.LevelFatal)

	randSeeder := rand.Reader
	key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, DefaultGenerateKeyPairBits, randSeeder)
	if err != nil {
		logrus.Errorf("unable to GenerateKeyPair for new peer: %v", err)
		key = emptyKey
	}
	logrus.Infof("New Peer: %s", n.String())
	runtimeID := uuid.New()
	logrus.Debugf("New Peer Runtime ID: %s", runtimeID.String())

	// Linux specific
	// This can fix the log line about UDP sizing
	//sysctl.Set("net.core.rmem_max", "2500000")
	// Linux specific

	return &Peer{
		Name:        n,
		uniqKey:     key,
		established: false,
	}
}

func Self() *Peer {
	return NewPeer(name.New(common.Self))
}

func (p *Peer) Host() host.Host {
	return p.host
}

func (p *Peer) ConnectPeerString(to string) error {
	peerID, err := peer.Decode(to)
	if err != nil {
		return fmt.Errorf("unable to decode string to peer id: %v", err)
	}
	return p.ConnectPeer(peerID)
}

func (p *Peer) ConnectPeer(to peer.ID) error {

	err := p.Handshake(to) // Not necessarily *required* but it's a good check for basic connectivity
	if err != nil {
		return fmt.Errorf("unable to initialize required handshake before grpc: %v", err)
	}

	logrus.Infof("gRPC Dial to: %s", to.String())
	grpcProto := p2pgrpc.NewGRPCProtocol(context.Background(), p.Host())
	conn, err := grpcProto.Dial(context.Background(), to, grpc.WithTimeout(time.Second*10), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	err = p.clientConnect(conn)
	if err != nil {
		return fmt.Errorf("unable to establish connection: %v", err)
	}

	return nil
}

func (p *Peer) ConnectSock(sock string) error {

	// Cache the socket
	p.localSocket = sock

	logrus.Warnf("mTLS disabled. running insecure.")
	conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", p.localSocket),
		grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return fmt.Errorf("failed grpc.Dial local socket: %v", err)
	}
	return p.clientConnect(conn)
}

func (p *Peer) clientConnect(conn grpc.ClientConnInterface) error {
	// Establish the connection from the conn
	core := rpc.NewCoreClient(conn)
	p.CoreClient = core
	runtime := rpc.NewRuntimeClient(conn)
	p.RuntimeClient = runtime
	schedule := rpc.NewScheduleClient(conn)
	p.ScheduleClient = schedule
	proxy := rpc.NewProxyClient(conn)
	p.ProxyClient = proxy
	logrus.Warnf("Connected to grpc")

	return nil
}

func (p *Peer) Establish(ctx context.Context, offset int) error {

	// [Host]
	//
	// Create a host with the Aurae default options
	basicHost, err := libp2p.New(DefaultOptions(p.uniqKey, offset)...)
	if err != nil {
		return err
	}

	// [DHT]
	//
	// Create a new distributed hash table for storing records
	dstore := dsync.MutexWrap(ds.NewMapDatastore())
	dht := dht.NewDHT(ctx, basicHost, dstore)

	// Routed Host
	routedHost := rhost.Wrap(basicHost, dht)
	p.host = routedHost
	p.established = true

	// Bootstrap
	err = p.Bootstrap(IPFSPeers)
	if err != nil {
		logrus.Errorf("Unable to bootstrap with IPFS: %v", err)
		p.established = false
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	// ID
	logrus.Infof("ID: %s", routedHost.ID().Pretty())
	return nil
}
