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
	"github.com/kris-nova/aurae/pkg/hostname"
	"testing"
)

func TestGraph1to1(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
	}
	graph := graph1to1()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

func TestGraph3Cycle(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
		2: "c",
	}
	graph := graph3cycle()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

func TestGraph5OuterCycle(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
	}
	graph := graph5cycleOuter()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

func TestGraph5SingleInnerLink(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
	}
	graph := graph5cycleSingleInnerLink()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

func TestGraph5FullInnerLink(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
	}
	graph := graph5cycleFullInnerLink()
	actual := CalculateHamiltonianPathHostname(graph)
	if !AssertHamPath(actual, expected) {
		t.Errorf("Unable to find Ham path. Actual: %v, Expected: %v", actual, expected)
	}
}

func TestGraphSuperNode(t *testing.T) {
	expected := TestHamPathCase{
		0: "a",
		1: "b",
		2: "c",
		3: "d",
		4: "e",
	}
	graph := graphSuperNode()
	actual := CalculateHamiltonianPathHostname(graph)
	if AssertHamPath(actual, expected) {
		t.Errorf("Expected failure. Actual: %+v, Expected: %v", actual, expected)
	}
}

type TestHamPathCase map[int]string

func AssertHamPath(a HamiltonianPath, t TestHamPathCase) bool {
	for _, peer := range a {
		found := false
		for _, name := range t {
			if name == peer.Hostname.Host {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// a ----- b
func graph1to1() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	b.AddPeer(root)
	return root
}

// a ----- b
//  \     /
//   \   /
//     c
func graph3cycle() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	c := b.ToPeer(hostname.New("c"))
	c.AddPeer(root)
	return root
}

// a ----- b
// |       |
// e       c
//  \     /
//     d
func graph5cycleOuter() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	c := b.ToPeer(hostname.New("c"))
	d := c.ToPeer(hostname.New("d"))
	e := d.ToPeer(hostname.New("e"))
	e.AddPeer(root)
	return root
}

// a ----- b
// |       |
// e ----- c
//  \     /
//     d
func graph5cycleSingleInnerLink() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	c := b.ToPeer(hostname.New("c"))
	d := c.ToPeer(hostname.New("d"))
	e := d.ToPeer(hostname.New("e"))

	c.AddPeer(e) // Single inner link

	e.AddPeer(root)
	return root
}

// a ----- b
// |   X   |
// e ----- c
//  \     /
//     d
func graph5cycleFullInnerLink() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	c := b.ToPeer(hostname.New("c"))
	d := c.ToPeer(hostname.New("d"))
	e := d.ToPeer(hostname.New("e"))

	c.AddPeer(e)
	e.AddPeer(c)

	root.AddPeer(c)
	root.AddPeer(e)
	b.AddPeer(e)
	e.AddPeer(b)

	e.AddPeer(root)
	return root
}

//   y
// / | \
// x---z
// \ | /
//   a ----- b
//   |   X   |
//   e ----- c
//    \     /
//       d
//
func graphSuperNode() *Peer {
	root := NewPeer("a", nil)
	b := root.ToPeer(hostname.New("b"))
	c := b.ToPeer(hostname.New("c"))
	d := c.ToPeer(hostname.New("d"))
	e := d.ToPeer(hostname.New("e"))

	c.AddPeer(e)
	e.AddPeer(c)

	root.AddPeer(c)
	root.AddPeer(e)
	b.AddPeer(e)
	e.AddPeer(b)

	e.AddPeer(root)

	// Create a sub graph that cycles
	x := root.ToPeer(hostname.New("x"))
	y := root.ToPeer(hostname.New("y"))
	z := root.ToPeer(hostname.New("z"))
	x.AddPeer(y)
	y.AddPeer(x)
	y.AddPeer(z)
	z.AddPeer(y)
	z.AddPeer(x)
	x.AddPeer(z)

	return root
}
