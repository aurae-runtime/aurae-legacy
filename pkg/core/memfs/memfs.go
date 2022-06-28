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
	"sync"
)

var rootNode = &Node{
	Name:     "/",
	Children: make(map[string]*Node),
	depth:    0,
}

var mtx = sync.Mutex{}

func Get(key string) string {
	mtx.Lock()
	defer mtx.Unlock()
	path := common.Path(key) // Data mutation!
	sub := rootNode.GetSubNode(path)
	if sub != nil {
		return string(sub.Content)
	}
	return ""
}

func Set(key, value string) {
	mtx.Lock()
	defer mtx.Unlock()
	path := common.Path(key) // Data mutation!
	rootNode.AddSubNode(path, value)
}

func Remove(key string) {
	mtx.Lock()
	defer mtx.Unlock()
	path := common.Path(key) // Data mutation!
	node := rootNode.GetSubNode(path)
	node.RemoveRecursive() // This will remove the node itself unless it is root!
}

func List(key string) map[string]bool {
	mtx.Lock()
	defer mtx.Unlock()
	base := common.Path(key)
	lsMap := rootNode.ListSubNodes(base)
	ret := make(map[string]bool)
	for file, node := range lsMap {
		ret[file] = node.File
	}
	return ret
}

func ListNode(key string) map[string]*Node {
	mtx.Lock()
	defer mtx.Unlock()
	base := common.Path(key)
	lsMap := rootNode.ListSubNodes(base)
	ret := make(map[string]*Node)
	for file, node := range lsMap {
		ret[file] = node
	}
	return ret
}
