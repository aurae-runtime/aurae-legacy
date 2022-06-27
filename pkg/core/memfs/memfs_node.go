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

const (
	DefaultNodePerm os.FileMode = 0644
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

	// parent is the parent of this *Node.
	//
	// parent is nil for the root of the tree.
	parent *Node

	// depth is the depth of this Node from the root.
	// The root Node counts as 1, so the "/file/path" Node would have a depth
	// of 3.
	depth int

	// file=true The Node is a file (with content)
	// file=false The Node is a directory (without content)
	File bool
}

// AddSubNode is how to introduce a new sub node.
//
// Directories are created recursively, and as needed.
// There is no concept of creating an empty directory.
func (n *Node) AddSubNode(key, value string) *Node {
	key = strings.TrimPrefix(key, "/")
	child := &Node{
		Children: make(map[string]*Node),
		depth:    n.depth + 1,
		File:     false,
		parent:   n,
		perm:     DefaultNodePerm,
	}
	if key == "" {
		return n
	}

	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		// See if the node already has the child
		if cachedChild, ok := n.Children[spl[0]]; ok {
			child = cachedChild
		}
		child.Name = spl[0]
		child.File = false
		child.AddSubNode(strings.Join(spl[1:], "/"), value)
	} else {
		child.Name = key
		child.File = true
		child.Content = []byte(value)
	}

	// Check if the node exists
	if _, ok := n.Children[child.Name]; ok {
		n.Children[child.Name].Content = []byte(value)
		return n.Children[child.Name]
	}

	// In some cases we silently add child nodes to a directory, so when
	// we add a child we also turn file=false
	n.File = false
	n.Children[child.Name] = child
	return child
}

// GetSubNode will return a sub Node recursively if it is found in the tree.
func (n *Node) GetSubNode(key string) *Node {
	key = strings.TrimPrefix(key, "/")
	if key == "" {
		return rootNode
	}
	if n.Name == key && n.File {
		return n
	}
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		first := spl[0]
		for _, child := range n.Children {
			if child.Name == first {
				return child.GetSubNode(strings.Join(spl[1:], "/"))
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

// ListSubNodes will return a flat listing of all children of a Node in the tree.
func (n *Node) ListSubNodes(key string) map[string]*Node {
	key = strings.TrimPrefix(key, "/")
	if key == "" || key == "/" {
		return rootNode.Children
	}
	result := make(map[string]*Node)

	// First check and see if a dir
	found := n.GetSubNode(key)
	if found == nil {
		return result // Nothing
	}
	if found.File {
		result[found.Name] = found
	}
	if !found.File {
		for _, c := range found.Children {
			if c.File {
				result[c.Name] = c
			} else {
				result[c.Name] = nil
			}
		}
	}
	return result
}

// RemoveRecursive will remove this node, and all of its children from the tree
func (n *Node) RemoveRecursive() {
	for _, c := range n.Children {
		c.RemoveRecursive()
	}
	n.remove()
}

// remove will remove this node from the tree. This is a dangerous operation
// to call externally, so we only export RemoveRecursive to other packages.
func (n *Node) remove() {
	if n.parent == nil {
		// Root node
		return
	}
	delete(n.parent.Children, n.Name)
}

func (n *Node) TotalChildren() int {
	var i int
	if n.parent == nil {
		i = 0
	} else {
		i = 1
	}
	for _, c := range n.Children {
		i = i + c.TotalChildren()
	}
	return i
}
