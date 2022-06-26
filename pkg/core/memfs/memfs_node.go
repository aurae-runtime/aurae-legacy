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

package memfs

import (
	"strings"
)

type Node struct {
	Name     string
	Value    string
	Children map[string]*Node
	depth    int
	file     bool
}

func (n *Node) AddChild(key, value string) *Node {
	key = strings.TrimSuffix(key, "/")
	child := &Node{
		Children: make(map[string]*Node),
		depth:    n.depth + 1,
		file:     false,
	}
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		// See if the node already has the child
		if cachedChild, ok := n.Children[spl[0]]; ok {
			child = cachedChild
		}
		child.Name = spl[0]
		child.file = false
		child.AddChild(strings.Join(spl[1:], "/"), value)
	} else {
		child.Name = key
		child.file = true
		child.Value = value
	}
	n.Children[child.Name] = child
	return child
}

func (n *Node) GetChild(key string) *Node {
	key = strings.TrimSuffix(key, "/")
	if n.Name == key && n.file {
		return n
	}
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		first := spl[0]
		for _, child := range n.Children {
			if child.Name == first {
				return child.GetChild(strings.Join(spl[1:], "/"))
			}
		}
	} else {
		for _, child := range n.Children {
			if child.Name == key {
				return child
			}
		}
	}
	return nil
}

func (n *Node) ListChildren(key string) map[string]*Node {
	result := make(map[string]*Node)
	key = strings.TrimSuffix(key, "/")
	// First check and see if its a dir
	found := rootNode.GetChild(key)
	if found == nil {
		return result // Nothing
	}
	if found.file {
		result[found.Name] = found
	}
	if !found.file {
		for _, c := range found.Children {
			if c.file {
				result[c.Name] = c
			} else {
				result[c.Name] = nil
			}
		}
	}
	return result
}
