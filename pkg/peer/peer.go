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
	"fmt"
	"github.com/google/uuid"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sirupsen/logrus"
)

const (
	AuraeStream              string = "/aurae"    // The official stream endpoint for Aurae
	AuraeStreamVersionFormat string = "/aurae/%s" // Format with the package version
)

func AuraeStreamProtocol() protocol.ID {
	auraeStreamProtocol := fmt.Sprintf(AuraeStreamVersionFormat, aurae.Version)
	ids := protocol.ConvertFromStrings([]string{auraeStreamProtocol})
	if len(ids) != 1 {
		panic("unable to find aurae protocol!")
	}
	return ids[0]
}

// Peer represents a single peer in the mesh.
//
// Peers exist in a directed graph.
//
// Each peer will be mapped by its Hostname in a local hash table.
// Each peer will be mapped by its Hostname in a public DHT.
// Each peer will dependent on a cryptographic key in order to initialize.
type Peer struct {

	// Name is the 3 part name of the peer in the mesh.
	Name *name.Name

	// Peers is where the digraph happens.
	Peers map[string]*Peer

	// Host is the peer instance of this peer.
	Host host.Host

	// DNS is an instance of multicast DNS
	DNS mdns.Service

	// runtimeID is a UUID generated at runtime
	// that exists for this specific reference
	// to the Peer in the network.
	//
	// This ID should never "persist" past the
	// execution context of this particular runtime.
	runtimeID uuid.UUID

	// All hosts are encrypted by default on the public
	Key crypto.PrivKey

	// established denotes if a peer is established in the mesh
	// or not
	established bool
}

// NewPeer will initialize a new *Peer without connecting.
//
// This will be an empty reference, and will do nothing
// until Connect() is called.
func NewPeer(n *name.Name, key crypto.PrivKey) *Peer {
	return &Peer{
		Peers:     make(map[string]*Peer),
		Name:      n,
		runtimeID: uuid.New(),
		Key:       key,
	}
}

func NewPeerServicename(svc string, key crypto.PrivKey) *Peer {
	return NewPeer(name.New(svc), key)
}

// Establish will initialize a network connection with the peer to peer circuit.
// Establish will also initialize multicast domain name name (mDNS) for
// managing distributed name names.
func (p *Peer) Establish() (host.Host, error) {

	// [p2p]
	// Here is where we establish ourselves in the mesh.
	h, err := p2p.New(DefaultOptions(p.Key)...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize peer-to-peer host: %v", err)
	}
	p.Host = h
	p.Host.SetStreamHandler(AuraeStreamProtocol(), func(s network.Stream) {
		logrus.Infof("Received stream: %v", s.ID())
	})
	logrus.Infof("Established. Listening on: %v", h.Network().ListenAddresses())

	// [mDNS]
	// Here is where we identify ourselves in the mesh.
	dns := mdns.NewMdnsService(h, p.Name.Service(), Notifee())
	p.DNS = dns
	logrus.Infof("Multicast DNS Established. Hostname: %s", p.Name.Service())

	return h, nil
}

func (p *Peer) ID() peer.ID {
	return p.Host.ID()
}

var self *Peer

// Self is a singleton for one's self in the mesh.
//
// TODO we need a way to cleanly manage service names
// TODO if the peer already exists.
func Self(key crypto.PrivKey) *Peer {
	if self == nil {
		self = NewPeer(name.New(common.Localhost), key)
	}
	return self
}

// --
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// --

// ToPeer will initialize a new Peer object based on a
// name.
//
// This mechanism will effectively serve as an alternative
// to DNS for the mesh if the peer is able to connect.
//
// Note: Connect() MUST be called on the new peer outside
// the scope of this function. This is effectively an AddChild()
// style function.
//func (p *Peer) ToPeer(h *name.Hostname) *Peer {
//	newPeer := NewPeerFromHostname(h, p.Key)
//	p.AddPeer(newPeer)
//	return newPeer
//}

//func (p *Peer) AddPeer(newPeer *Peer) {
//	p.Peers[newPeer.Name.String()] = newPeer
//}
