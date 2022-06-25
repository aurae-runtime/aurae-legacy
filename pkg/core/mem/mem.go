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

import "sync"

type Database struct {
	mtx sync.Mutex
}

func NewDatabase() *Database {
	return &Database{
		mtx: sync.Mutex{},
	}
}

var packageCache = map[string]string{}

func (c *Database) Get(key string) string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if r, ok := packageCache[key]; ok {
		return r
	}
	return ""
}

func (c *Database) Set(key, value string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	packageCache[key] = value
}

func (c *Database) List(key string) map[string]string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	// TODO implement filesystem query by paths
	return packageCache
}
