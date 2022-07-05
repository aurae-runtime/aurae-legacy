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
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
)

// Establish will join the mesh.
//
//
//
func (p *Peer) Establish() error {

	ctx := context.Background() // TODO Pass this in or generate this somehow

	// [p2p]
	// Here is where we establish ourselves in the mesh.
	h, err := p2p.New(DefaultOptions()...)
	if err != nil {
		return fmt.Errorf("unable to initialize peer-to-peer host: %v", err)
	}
	p.Host = h
	p.Host.SetStreamHandler(AuraeStreamProtocol(), func(s network.Stream) {
		logrus.Infof("Received stream: %v", s.ID())
	})
	logrus.Infof("Established. Listening on: %v", h.Network().ListenAddresses())

	// [Bootstrap]
	dstore := dsync.MutexWrap(ds.NewMapDatastore())
	dht := dht.NewDHT(ctx, h, dstore)
	routedHost := rhost.Wrap(h, dht)
	p.RHost = routedHost
	logrus.Infof("Bootstrapping with %d peers.", len(IPFSPeers))
	err = p.Bootstrap(IPFSPeers)
	if err != nil {
		return fmt.Errorf("unable to bootstrap routed host with IPFS peers: %v", err)
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		return fmt.Errorf("unable to bootstrap DHT: %v", err)
	}

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

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))
	for _, addr := range p.RHost.Addrs() {
		logrus.Infof("Listening on address: %s", addr.Encapsulate(hostAddr))
	}

	logrus.Infof("Public address to connect: %s", p.RHost.ID().Pretty())

	logrus.Infof("Established.")
	return nil
}
