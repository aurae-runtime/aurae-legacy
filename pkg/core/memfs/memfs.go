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
	"github.com/kris-nova/aurae/pkg/common"
	"strings"
	"sync"
)

type Database struct {
	mtx sync.Mutex
}

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

func (n *Node) ListChildren(key string) map[string]string {
	result := make(map[string]string)
	key = strings.TrimSuffix(key, "/")
	// First check and see if its a dir
	found := rootNode.GetChild(key)
	if found == nil {
		return result // Nothing
	}
	if found.file {
		result[found.Name] = found.Value
	}
	if !found.file {
		for _, c := range found.Children {
			if c.file {
				result[c.Name] = found.Value
			} else {
				result[c.Name] = ""
			}
		}
	}
	return result
}

func NewDatabase() *Database {
	return &Database{
		mtx: sync.Mutex{},
	}
}

var rootNode = &Node{
	Name:     "/",
	Children: make(map[string]*Node),
	depth:    0,
}

func (c *Database) Get(key string) string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	path := common.Path(key) // Data mutation!
	return rootNode.GetChild(path).Value
}

func (c *Database) Set(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	path := common.Path(key) // Data mutation!
	rootNode.AddChild(path, value)
}

func (c *Database) List(key string) map[string]string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	base := common.Path(key)
	return rootNode.ListChildren(base)
}
