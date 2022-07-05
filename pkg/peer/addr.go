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
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// Address will calculate an address encapsulated from IPFS
//
// Example address: /ip4/127.0.0.1/tcp/10000/p2p/QmYo41GybvrXk8y8Xnm1P7pfA4YEXCpfnLyzgRPnNbG35e
func (p *Peer) Address() string {
	// Create a new IPFS specific address
	ipfsAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", p.Host.ID().Pretty()))
	if err != nil {
		logrus.Warnf("unable to create IPFS host address: %v", err)
		return ""
	}

	hostAddrs := p.Host.Addrs()
	if len(hostAddrs) < 1 {
		logrus.Warnf("unable to find host address")
		return ""
	}
	hostAddr := hostAddrs[0]
	return hostAddr.Encapsulate(ipfsAddr).String()
}

// AddressDecode will decode a raw address from p.Address() and return
// the details needed to connect from another peer.
//
// This is taken from the sample code, and doesn't need to be this offensive.
//
// Pass the value that p.Address() returns here.
func AddressDecode(addr string) (peer.ID, multiaddr.Multiaddr) {

	ipfsAddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		logrus.Warnf("unable to create IPFS host address: %v", err)
		return "", ipfsAddr
	}

	ipfsProtocol, err := ipfsAddr.ValueForProtocol(multiaddr.P_IPFS)
	if err != nil {
		logrus.Warnf("unable to lookup ipfs protocol: %v", err)
		return "", ipfsAddr
	}

	peerID, err := peer.Decode(ipfsProtocol)
	if err != nil {
		logrus.Warnf("unable to decode ipfs protocol: %v", err)
		return "", ipfsAddr
	}

	// Here be dragons
	// TODO we can debug this and simplify this, it doesnt need to be this bad
	targetPeerAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ipfsProtocol))
	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)

	return peerID, targetAddr
}
