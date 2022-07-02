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

// HamiltonianPaths will be 0 indexed

type HamiltonianPathPeer map[int]*Peer
type HamiltonianPathHostname map[int]string

func NewHamiltonianPathPeer() HamiltonianPathPeer {
	return make(map[int]*Peer)
}

func NewHamiltonianPathHostname() HamiltonianPathHostname {
	return make(map[int]string)
}

func CalculateHamiltonianPathHostname(root *Peer) HamiltonianPathHostname {
	x := CalculateHamiltonianPathPeer(root)
	y := NewHamiltonianPathHostname()
	for i, peer := range x {
		y[i] = peer.Hostname
	}
	return y
}

func CalculateHamiltonianPathPeer(root *Peer) HamiltonianPathPeer {
	x := HamiltonianPathPeer{
		0: root,
	}
	return x
}
