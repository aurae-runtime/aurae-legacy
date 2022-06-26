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
	return string(rootNode.GetChild(path).Content)
}

func Set(key, value string) {
	mtx.Lock()
	defer mtx.Unlock()
	path := common.Path(key) // Data mutation!
	rootNode.AddChild(path, value)
}

func List(key string) map[string]string {
	mtx.Lock()
	defer mtx.Unlock()
	base := common.Path(key)
	lsMap := rootNode.ListChildren(base)
	ret := make(map[string]string)
	for file, node := range lsMap {
		if node == nil {
			ret[file] = ""
		} else {
			ret[file] = string(node.Content)
		}
	}
	return ret
}
