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
	"crypto/rand"
	p2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

// DefaultOptions are the default p2p options for Aurae
func DefaultOptions() []p2p.Option {
	priv, _, _ := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	return []p2p.Option{

		p2p.Identity(priv),

		// EnableAutoRelay for performance reasons in the circuit.
		p2p.EnableAutoRelay(),

		//p2p.ListenAddrs(
		// Listen on ipv4, choose an available port
		//multiaddr.StringCast("/ip4/0.0.0.0/tcp/0")),

		// Listen on the default p2p-circuit
		//multiaddr.StringCast("/p2p-circuit"),

		// Listen on IPv6, choose an available port
		//multiaddr.StringCast("/ip6/::/tcp/0")),
	}

}
