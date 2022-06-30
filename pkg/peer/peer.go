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
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/common"
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
	Hostname string
	Peers    map[string]*Peer
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
func (s *Peer) ToPeer(hostname string) *Peer {
	newPeer := &Peer{
		Hostname: hostname,
	}
	s.Peers[hostname] = newPeer
	return newPeer
}

// Connect will connect to the peer
func (*Peer) Connect() (*client.Client, error) {

	return nil, nil
}

// NewSafeConnection will return a new net.Conn
// from the Go standard library for the new peer.
//
// These connections MUST be safe to use while adhering
// the scope of the Aurae project.
func (*Peer) NewSafeConnection() *net.Conn {

	return nil
}
