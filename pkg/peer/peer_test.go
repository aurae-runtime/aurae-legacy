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

// An Euler path is a path that passes through every edge exactly once.
// If it ends at the initial vertex then it is an Euler cycle.
// A Hamiltonian path is a path that passes through every vertex exactly once (NOT every edge).
// If it ends at the initial vertex then it is a Hamiltonian cycle.

package peer

import (
	"testing"
)

func TestGraph(t *testing.T) {
	expected := AssertHamPathHostname{
		0: "a",
		1: "b",
	}
	graph := graph1to1()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

type AssertHamPathHostname map[int]string

func AssertHamPath(h HamiltonianPathHostname, a AssertHamPathHostname) bool {
	for i, hostname := range a {
		if x, ok := h[i]; ok {
			if x == hostname {
				continue
			}
		}
		return false
	}
	return true
}

// a ----- b
func graph1to1() *Peer {
	root := NewPeer("a")
	root.ToPeer("b")
	return root
}

// a ----- b
//  \     /
//   \   /
//     c
func graph3cycle() *Peer {
	root := NewPeer("a")
	b := root.ToPeer("b")
	c := b.ToPeer("c")
	c.AddPeer(root)
	return root
}

// 	peerNames := []string{"a", "b", "c", "1", "2", "3"}
//	for _, hostname := range peerNames {
//		root.ToPeer(hostname)
//	}
