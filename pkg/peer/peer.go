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
	"crypto"
	"fmt"
	"github.com/google/uuid"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/hostname"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"net"
)

// Peer represents a single peer in the mesh.
//
// Peers exist in a directed graph.
//
// Each peer will be mapped by its Hostname in a local hash table.
// Each peer will be mapped by its Hostname in a public DHT.
// Each peer will dependent on a cryptographic key in order to initialize.
type Peer struct {

	// Hostname is the 3 part name of the peer in the mesh.
	Hostname *hostname.Hostname

	// Peers is where the digraph happens.
	Peers map[string]*Peer

	// Host is the peer instance of this peer.
	Host host.Host

	// runtimeID is a UUID generated at runtime
	// that exists for this specific reference
	// to the Peer in the network.
	//
	// This ID should never "persist" past the
	// execution context of this particular runtime.
	runtimeID uuid.UUID

	// All hosts are encrypted by default on the public
	Key crypto.PrivateKey
}

var self *Peer

// Self is a singleton for one's self in the mesh.
func Self() *Peer {
	if self == nil {
		self = &Peer{
			Hostname: hostname.New(common.Self),
		}
	}
	return self
}

// ToPeer will initialize a new Peer object based on a
// hostname.
//
// This mechanism will effectively serve as an alternative
// to DNS for the mesh if the peer is able to connect.
//
// Note: Connect() MUST be called on the new peer outside
// the scope of this function. This is effectively an AddChild()
// style function.
func (p *Peer) ToPeer(h *hostname.Hostname) *Peer {
	newPeer := NewPeerFromHostname(h, p.Key)
	p.AddPeer(newPeer)
	return newPeer
}

// Connect will connect this peer to the mesh, and begin hosting
// this peer. After a Connect() is successful this peer can now
// accept connections from clients.
func (p *Peer) Connect() (host.Host, error) {
	host, err := p2p.New(DefaultOptions()...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize peer-to-peer host: %v", err)
	}
	p.Host = host
	p.Host.SetStreamHandler("/aurae", func(s network.Stream) {
		logrus.Infof("Received stream: %v", s.ID())
	})
	return host, nil
}

// NewSafeConnection will return a new net.Conn
// from the Go standard library for the new peer.
//
// These connections MUST be safe to use while adhering
// the scope of the Aurae project.
func (p *Peer) NewSafeConnection() *net.Conn {

	return nil
}

func (p *Peer) AddPeer(newPeer *Peer) *net.Conn {
	p.Peers[newPeer.Hostname.String()] = newPeer
	return nil
}

func (p *Peer) GetID() string {
	if p.Host == nil {
		return ""
	}
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", p.Host.ID().Pretty()))
	addr := p.Host.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}

func (p *Peer) DialID(id string) error {
	hostAddr, err := multiaddr.NewMultiaddr(id)
	if err != nil {
		return fmt.Errorf("unable to calculate multi address: %v", err)
	}
	pid, err := hostAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		return fmt.Errorf("unable to calculate protocol id: %v", err)
	}
	peerID, err := peer.Decode(pid)
	if err != nil {
		return fmt.Errorf("unable to decode peer id: %v", err)
	}
	targetPeerAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
	targetAddr := hostAddr.Decapsulate(targetPeerAddr)

	// TODO Understand the peer store more
	logrus.Infof("Adding to Peer store. PeerID: %s, Target Addr: %s,", peerID, targetAddr)
	p.Host.Peerstore().AddAddr(peerID, targetAddr, peerstore.PermanentAddrTTL)

	stream, err := p.Host.NewStream(context.Background(), peerID, "/aurae")
	if err != nil {
		return fmt.Errorf("unable to connect to stream: %v", err)
	}

	stream.Conn() // TODO manage net conn :)

	return nil
}

// NewPeer will initialize a new *Peer without connecting.
//
// This will be an empty reference, and will do nothing
// until Connect() is called.
func NewPeer(name string, key crypto.PrivateKey) *Peer {
	return NewPeerFromHostname(hostname.New(name), key)
}

// NewPeerFromHostname will initialize a new peer directly from a hostname.Hostname
func NewPeerFromHostname(hn *hostname.Hostname, key crypto.PrivateKey) *Peer {
	return &Peer{
		Peers:     make(map[string]*Peer),
		Hostname:  hn,
		runtimeID: uuid.New(),
		Key:       key,
	}
}
