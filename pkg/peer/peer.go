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
	"github.com/kris-nova/aurae/pkg/hostname"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
)

const (
	AuraeStream protocol.ID = "/aurae" // The official stream endpoint for Aurae
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

	// peerID is the ID used to index on in the distributed hash table.
	peerID peer.ID

	// peerAddr is this unique address in the mesh
	peerAddr multiaddr.Multiaddr

	// runtimeID is a UUID generated at runtime
	// that exists for this specific reference
	// to the Peer in the network.
	//
	// This ID should never "persist" past the
	// execution context of this particular runtime.
	runtimeID uuid.UUID

	// All hosts are encrypted by default on the public
	Key crypto.PrivKey
}

var self *Peer

// Self is a singleton for one's self in the mesh.
func Self(key crypto.PrivKey) *Peer {
	if self == nil {
		self = NewPeer(common.Localhost, key)
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
	h, err := p2p.New(DefaultOptions(p.Key)...)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize peer-to-peer host: %v", err)
	}
	p.Host = h
	p.Host.SetStreamHandler(AuraeStream, func(s network.Stream) {
		logrus.Infof("Received stream: %v", s.ID())
	})
	p.GetID()
	defer p.GetAddr()
	return h, nil
}

// NewSafeConnection will return a new net.Conn
// from the Go standard library for the new peer.
//
// These connections MUST be safe to use while adhering
// the scope of the Aurae project.
func (p *Peer) NewSafeConnection() (*net.Conn, error) {
	stream, err := p.Host.NewStream(context.Background(), p.peerID, AuraeStream)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to stream: %v", err)
	}
	if conn, ok := stream.Conn().(net.Conn); ok {
		return &conn, nil
	}
	return nil, fmt.Errorf("unable to convert to *net.Conn")
}

func (p *Peer) AddPeer(newPeer *Peer) {
	p.Peers[newPeer.Hostname.String()] = newPeer
}

func (p *Peer) GetID() string {
	if p.Host == nil {
		return ""
	}
	if p.peerID.String() == "" {
		hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", p.Host.ID().Pretty()))
		addr := p.Host.Addrs()[0]
		peerIDStr := addr.Encapsulate(hostAddr).String()
		peerID, err := peer.Decode(peerIDStr)
		if err != nil {
			return ""
		}
		p.peerID = peerID
	}
	return p.peerID.String()
}

func (p *Peer) GetAddr() (multiaddr.Multiaddr, error) {
	if p.peerAddr == nil {
		hostAddr, err := multiaddr.NewMultiaddr(p.peerID.String())
		if err != nil {
			return nil, fmt.Errorf("unable to calculate multi address: %v", err)
		}
		pid, err := hostAddr.ValueForProtocol(multiaddr.P_IPFS)
		if err != nil {
			return nil, fmt.Errorf("unable to calculate protocol id: %v", err)
		}
		peerID, err := peer.Decode(pid)
		if err != nil {
			return nil, fmt.Errorf("unable to decode peer id: %v", err)
		}
		targetPeerAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", pid))
		targetAddr := hostAddr.Decapsulate(targetPeerAddr)

		logrus.Infof("Adding to Peer store. PeerID: %s, Target Addr: %s,", peerID, targetAddr)
		p.Host.Peerstore().AddAddr(peerID, targetAddr, peerstore.PermanentAddrTTL)
		p.peerAddr = targetAddr
	}
	return p.peerAddr, nil
}

// NewPeer will initialize a new *Peer without connecting.
//
// This will be an empty reference, and will do nothing
// until Connect() is called.
func NewPeer(name string, key crypto.PrivKey) *Peer {
	return NewPeerFromHostname(hostname.New(name), key)
}

// NewPeerFromHostname will initialize a new peer directly from a hostname.Hostname
func NewPeerFromHostname(hn *hostname.Hostname, key crypto.PrivKey) *Peer {
	return &Peer{
		Peers:     make(map[string]*Peer),
		Hostname:  hn,
		runtimeID: uuid.New(),
		Key:       key,
	}
}

func KeyFromPath(path string) (crypto.PrivKey, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// TODO We should check the keys and support all the SSH keys
	return crypto.UnmarshalEd25519PrivateKey(bytes)
	//return crypto.UnmarshalPrivateKey(bytes)
}
