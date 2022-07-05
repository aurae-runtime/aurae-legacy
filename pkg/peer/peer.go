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
	"fmt"
	"github.com/google/uuid"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peerstore"
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
	Name name.Name

	// Peers is where the digraph happens.
	Peers map[string]*Peer

	// Host is the peer instance of this peer.
	Host host.Host

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
	runtimeID uuid.UUID

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
		runtimeID: uuid.New(),
	}
}

func NewPeerServicename(svc string) *Peer {
	return NewPeer(name.New(svc))
}

// Establish will initialize a network connection with the peer to peer circuit.
// Establish will also initialize multicast domain name name (mDNS) for
// managing distributed name names.
func (p *Peer) Establish() (host.Host, error) {

	// [p2p]
	// Here is where we establish ourselves in the mesh.

	h, err := p2p.New(DefaultOptions()...)
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
	dns := NewNameService()
	internalDNS := mdns.NewMdnsService(h, p.Name.Service(), dns)
	internalDNS.Start()
	p.DNS = dns
	p.internalDNS = internalDNS
	logrus.Infof("Multicast DNS Established. Hostname: %s", p.Name.Service())

	// All Peers will respond on the default Aurae Protocol.
	// We establish that handler now.
	addr := p.Address()
	logrus.Infof("Establish Aurae protocol [%s] on address: %s", AuraeStreamProtocol(), addr)
	return h, nil
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
