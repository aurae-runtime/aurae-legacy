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
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	golog "github.com/ipfs/go-log/v2"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
)

const (
	DefaultGenerateKeyPairBits int = 2048
	DefaultListenPort          int = 8709
	DefaultPeerPort            int = 8708
)

var emptyKey crypto.PrivKey = &crypto.Ed25519PrivateKey{}

type Peer struct {
	uniqKey     crypto.PrivKey
	Name        name.Name
	Peers       map[string]*Peer
	Host        host.Host
	RHost       rhost.RoutedHost
	DNS         *NameService
	internalDNS mdns.Service
	RuntimeID   uuid.UUID
	established bool
}

func NewPeer(n name.Name) *Peer {
	golog.SetAllLoggers(golog.LevelPanic)
	golog.SetupLogging(golog.Config{
		Stdout: false,
		Stderr: false,
	})
	randSeeder := rand.Reader
	key, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, DefaultGenerateKeyPairBits, randSeeder)
	if err != nil {
		logrus.Errorf("unable to GenerateKeyPair for new peer: %v", err)
		key = emptyKey
	}
	logrus.Infof("New Peer: %s", n.String())
	return &Peer{
		Name:        n,
		uniqKey:     key,
		established: false,
	}
}

func Self() *Peer {
	return NewPeer(name.New(common.Self))
}

func (p *Peer) Establish(ctx context.Context, offset int) error {

	// Options
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", DefaultListenPort+offset)),
		libp2p.Identity(p.uniqKey),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}

	// Host
	basicHost, err := libp2p.New(opts...)
	if err != nil {
		return err
	}
	p.Host = basicHost

	// DHT
	dstore := dsync.MutexWrap(ds.NewMapDatastore())
	dht := dht.NewDHT(ctx, basicHost, dstore)

	// Routed Host
	routedHost := rhost.Wrap(basicHost, dht)
	p.RHost = *routedHost

	// Bootstrap
	err = p.Bootstrap(IPFSPeers)
	if err != nil {
		return err
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	// ID
	logrus.Infof("ID: %s", routedHost.ID().Pretty())
	p.established = true
	return nil
}

func (p *Peer) To(peerID string) error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.RHost.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		if err := doEcho(s); err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	id, err := peer.Decode(peerID)
	if err != nil {
		return err
	}

	s, err := p.RHost.NewStream(context.Background(), id, "/echo/1.0.0")
	if err != nil {
		return err
	}

	_, err = s.Write([]byte("Hello, world!\n"))
	if err != nil {
		return err
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}

	logrus.Infof("%q", out)
	return nil
}

func (p *Peer) Stream() error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.RHost.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		logrus.Infof("Received new stream: %s", s.ID())
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
	select {} // hang forever
}

func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	logrus.Infof("Read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}
