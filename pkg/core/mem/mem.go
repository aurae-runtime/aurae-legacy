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

package mem

import (
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
}

func (n *Node) AddChild(key, value string) *Node {
	child := &Node{
		Value:    value,
		Children: make(map[string]*Node),
		depth:    n.depth + 1,
	}
	spl := strings.Split(key, "/")
	if len(spl) > 1 {
		child.Name = spl[1]
		child.AddChild(strings.Join(spl[1:], "/"), value)
	} else {
		child.Name = key
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
		first := spl[1]
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
	} else if key == "*" || len(spl) == 1 {
		for _, child := range n.Children {
			result[child.Name] = ""
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
	return rootNode.GetChild(key).Value
}

func (c *Database) Set(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	rootNode.AddChild(key, value)
}

func (c *Database) List(key string) map[string]string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return rootNode.ListChildren(key)
}
