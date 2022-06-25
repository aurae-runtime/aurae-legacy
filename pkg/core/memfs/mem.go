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
	child := &Node{
		Value:    value,
		Children: make(map[string]*Node),
		depth:    n.depth + 1,
		file:     false,
	}
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		child.Name = spl[0]
		child.AddChild(strings.Join(spl[1:], "/"), value)
	} else {
		child.Name = key
		child.file = true
	}
	n.Children[child.Name] = child
	return child
}

func (n *Node) GetChild(key string) *Node {
	if n.Name == key {
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
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		first := spl[1]
		for _, child := range n.Children {
			if child.Name == first {
				return child.ListChildren(strings.Join(spl[1:], "/"))
			}
		}
	} else if key == "*" {
		for _, child := range n.Children {
			if child.file {
				result[child.Name] = child.Value // File
			} else {
				result[child.Name] = "" // Dir
			}
		}
	} else if len(spl) == 1 {
		for _, child := range n.Children {
			if child.Name == spl[0] {
				return child.ListChildren("*")
			}
		}
	} else {
		for _, child := range n.Children {
			if child.Name == key {
				return child.ListChildren(key)
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
	path := common.Path(key) // Data mutation!
	return rootNode.ListChildren(path)
}
