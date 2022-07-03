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

// Notes:
//
// We must be able to calculate the size of the graph dynamically.
// We will never be able to determine the size of the graph before
// walking the graph.

package peer

// HamiltonianPaths will be 0 indexed

type HamiltonianPath map[int]*Peer

func NewHamiltonianPath() HamiltonianPath {
	return make(map[int]*Peer)
}

// CalculateHamiltonianPath is where the magic happens.
//
// Here is where the magic happens.
func CalculateHamiltonianPath(root *Peer) HamiltonianPath {
	x := NewHamiltonianPath()
	if x.recursiveCycle(root, 0) {
		// We have found a path
		return x
	}
	x = NewHamiltonianPath() // Reset if no cycle is found (return empty)
	return x
}

// recursiveCycle will assert a single root against a graph
//
// This is where the core logic of the Hamilton path algorithm lives.
func (h HamiltonianPath) recursiveCycle(graph *Peer, pos int) bool {

	// Check if this peer already exists in the graph
	for _, x := range h {
		if x.runtimeID == graph.runtimeID {
			return true
		}
	}

	// Base case. Set the current position to the current *Peer in the graph.
	h[pos] = graph

	// All paths must cycle back to 0 in order for the path
	// to be a true Ham path.
	connectsToRoot := false
	for _, peer := range graph.Peers {
		if peer.runtimeID == h[0].runtimeID {
			connectsToRoot = true
		} else {
			return h.recursiveCycle(peer, pos+1)
		}
	}
	return connectsToRoot
}
