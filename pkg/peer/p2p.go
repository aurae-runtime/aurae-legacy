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
	"github.com/libp2p/go-libp2p-crypto"
	"github.com/multiformats/go-multiaddr"
)

// DefaultOptions are the default p2p options for Aurae
func DefaultOptions(key crypto.PrivKey) []p2p.Option {

	return []p2p.Option{

		// We MUST use a key in order to adhere to our guarantee that all
		// connections can be fundamentally achieved with only a single
		// encryption key.
		p2p.Identity(key),

		// EnableAutoRelay for performance reasons in the circuit.
		p2p.EnableAutoRelay(),

		p2p.ListenAddrs(
			// Listen on ipv4, choose an available port
			//multiaddr.StringCast("/ip4/0.0.0.0/tcp/0")),

			// Listen on the default p2p-circuit
			//multiaddr.StringCast("/p2p-circuit"),

			// Listen on IPv6, choose an available port
			multiaddr.StringCast("/ip6/::/tcp/0")),
	}

}
