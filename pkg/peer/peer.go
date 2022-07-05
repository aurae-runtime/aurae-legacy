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
	"github.com/google/uuid"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

// Peer represents a single peer in the mesh.
//
// Peers exist in a directed graph.
//
// Each peer will be mapped by its Hostname in a local hash table.
// Each peer will be mapped by its Hostname in a public DHT.
// Each peer will dependent on a cryptographic key in order to initialize.
type Peer struct {

	// Name is the 3 part name of the peer in the mesh.
	Name name.Name

	// Peers is where the digraph happens.
	Peers map[string]*Peer

	// Host is the peer instance of this peer.
	Host host.Host

	// RHost is a wrapper around Host
	RHost *routedhost.RoutedHost

	// DNS is an instance of multicast DNS
	DNS *NameService

	// internalDNS is the mdns service itself
	internalDNS mdns.Service

	// runtimeID is a UUID generated at runtime
	// that exists for this specific reference
	// to the Peer in the network.
	//
	// This ID should never "persist" past the
	// execution context of this particular runtime.
	RuntimeID uuid.UUID

	// established denotes if a peer is established in the mesh
	// or not
	established bool
}

// NewPeer will initialize a new *Peer without connecting.
//
// This will be an empty reference, and will do nothing
// until Connect() is called.
func NewPeer(n name.Name) *Peer {
	return &Peer{
		Peers:     make(map[string]*Peer),
		Name:      n,
		RuntimeID: uuid.New(),
	}
}

func NewPeerServicename(svc string) *Peer {
	return NewPeer(name.New(svc))
}

// ToPeer is a simple heartbeat method to exercise connectivity.
//func (p *Peer) ToPeer(n name.Name) error {
//	return fmt.Errorf("UNSUPPORTED")
//}

//func (p *Peer) ToPeerID(id peer.ID) (network.Stream, error) {
//
//}

// ToPeerAddr will dial an address directly.
//
// You can find an address by calling p.Address()
func (p *Peer) ToPeerAddr(addr string) (network.Stream, error) {
	id, ma := AddressDecode(addr)
	p.Host.Peerstore().AddAddr(id, ma, peerstore.PermanentAddrTTL)
	return p.Host.NewStream(context.Background(), id, AuraeStreamProtocol())
}

func (p *Peer) Close() error {
	return p.internalDNS.Close()
}

var self *Peer

// Self is a singleton for one's self in the mesh.
func Self() *Peer {
	if self == nil {
		self = NewPeer(name.New(common.Localhost))
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
