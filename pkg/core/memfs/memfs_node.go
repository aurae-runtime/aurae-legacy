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
	"os"
	"strings"
)

// Node represents the vertex in a Graph, or, an Inode in a file system.
//
// A Node can represent a directory, or a file.
//
// If a Node is a file (and not a directory) it's Node.file=true flag will
// be set.
type Node struct {

	// Name is the exact name of the Node, also known as its filename.
	// The Name should be unique for its subtree, in the same way that a
	// directory cannot have multiple files with the same name.
	Name string

	// Content is the content of the virtual file.
	Content []byte

	// perm is the permissions bits of the file to represent to the user.
	//
	// We carefully use the word "perm" here as "mode" has special meaning
	// in libfuse.
	perm os.FileMode

	// Sub nodes, or nested files and directories. Each with a unique name.
	Children map[string]*Node

	// depth is the depth of this Node from the root.
	// The root Node counts as 1, so the "/file/path" Node would have a depth
	// of 3.
	depth int

	// file=true The Node is a file (with content)
	// file=false The Node is a directory (without content)
	file bool
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
		child.Content = []byte(value)
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
