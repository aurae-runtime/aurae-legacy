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
	"crypto/rand"
	"github.com/google/uuid"
	golog "github.com/ipfs/go-log/v2"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sirupsen/logrus"
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
	host        host.Host
	dns         *NameService
	internalDNS mdns.Service
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
	logrus.Debugf("New Peer: %s", n.String())
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
