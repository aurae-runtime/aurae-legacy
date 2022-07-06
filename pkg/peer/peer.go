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
	uniqKey crypto.PrivKey
	Name    name.Name
	//Peers       map[string]*Peer
	Host  host.Host
	RHost rhost.RoutedHost
	//DNS         *NameService
	//internalDNS mdns.Service
	RuntimeID   uuid.UUID
	established bool
}

func NewPeer(n name.Name) *Peer {
	golog.SetupLogging(golog.Config{
		Stdout: false,
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
		RuntimeID:   runtimeID,
	}
}

func Self() *Peer {
	return NewPeer(name.New(common.Self))
}

func (p *Peer) Establish(ctx context.Context, offset int) error {

	// [Host]
	//
	// Create a host with the Aurae default options
	basicHost, err := libp2p.New(DefaultOptions(p.uniqKey, offset)...)
	if err != nil {
		return err
	}
	p.Host = basicHost

	// [DHT]
	//
	// Create a new distributed hash table for storing records
	dstore := dsync.MutexWrap(ds.NewMapDatastore())
	dht := dht.NewDHT(ctx, basicHost, dstore)

	// Routed Host
	routedHost := rhost.Wrap(basicHost, dht)
	p.RHost = *routedHost
	p.established = true

	// Bootstrap
	go func() {
		err = p.Bootstrap(IPFSPeers)
		if err != nil {
			logrus.Errorf("Unable to bootstrap with IPFS: %v", err)
			p.established = false
		}
	}()
	err = dht.Bootstrap(ctx)
	if err != nil {
		return err
	}

	// ID
	logrus.Infof("ID: %s", routedHost.ID().Pretty())
	return nil
}

func (p *Peer) To(peerID string) error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.RHost.SetStreamHandler(AuraeStreamProtocol(), func(s network.Stream) {
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

	s, err := p.RHost.NewStream(context.Background(), id, AuraeStreamProtocol())
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
	p.RHost.SetStreamHandler(AuraeStreamProtocol(), func(s network.Stream) {
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
