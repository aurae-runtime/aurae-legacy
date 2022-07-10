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
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/multiformats/go-multiaddr"
)

// DefaultOptions are the default p2p options for Aurae
//
// key is a crypto.PrivKey
// offset is a port offset
//   - Server [+1]
//   - Client [+0]
func DefaultOptions(key crypto.PrivKey, offset int) []p2p.Option {

	// We believe this is needed to generate a valid host ID
	return []p2p.Option{

		p2p.Identity(key),

		// EnableAutoRelay for performance reasons in the circuit.
		// This is also suggested for the NAT hole punching.
		p2p.EnableAutoRelay(),

		p2p.DefaultTransports,
		p2p.DefaultMuxers,
		p2p.DefaultSecurity,
		p2p.EnableNATService(),   // Convenience service for peers
		p2p.NATPortMap(),         // Try to open a port in the network firewall using UPnP
		p2p.EnableHolePunching(), // STUN hole punching

		p2p.ListenAddrs(
			// Listen on ipv4, choose an available port
			multiaddr.StringCast("/ip4/0.0.0.0/tcp/0"),

			// Listen on the default p2p-circuit
			multiaddr.StringCast("/p2p-circuit"),

			// Listen on IPv6, choose an available port
			multiaddr.StringCast("/ip6/::/tcp/0")),

		//p2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", DefaultListenPort+offset)),

	}

}
