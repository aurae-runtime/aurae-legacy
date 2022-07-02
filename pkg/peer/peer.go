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
	"github.com/kris-nova/aurae/pkg/common"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"net"
)

// Research
//
// Should we build a DAG?
// How do we handle "cluster wide" updates?
// What are the advantages of a DAG?
// What other types of graphs should we consider?
// Is there prior art here?

// Peer represents a single peer in the mesh.
//
// Each peer will be mapped by its Hostname and will serve
// as
type Peer struct {

	// Hostname is a special string. This is the actual
	// DNS name of this peer in the network.
	Hostname string

	// Peers is where the digraph happens.
	Peers map[string]*Peer

	Host host.Host

	// runtimeID is a UUID generated at runtime
	// that exists for this specific reference
	// to the Peer in the network.
	//
	// This ID should never "persist" past the
	// execution context of this particular runtime.
	runtimeID uuid.UUID
}

var self *Peer

// Self is a singleton for one's self in the mesh.
func Self() *Peer {
	if self == nil {
		self = &Peer{
			Hostname: common.Self,
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
func (p *Peer) ToPeer(hostname string) *Peer {
	newPeer := NewPeer(hostname)
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
	p.Peers[newPeer.Hostname] = newPeer
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
func NewPeer(hostname string) *Peer {
	return &Peer{
		Peers:     make(map[string]*Peer),
		Hostname:  hostname,
		runtimeID: uuid.New(),
	}
}
