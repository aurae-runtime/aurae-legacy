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

package core

import (
	"context"
	"github.com/kris-nova/aurae/pkg/core/mem"
	"github.com/kris-nova/aurae/rpc"
	"sync"
)

var _ rpc.CoreServiceServer = &Database{}

const (
	CoreCode_OKAY int32 = 0
	CoreCode_ERR  int32 = 1
)

type Database struct {
	mtx sync.Mutex
	rpc.UnimplementedCoreServiceServer
	// TODO we will need to prioritize our getters and setters
	// TODO we need to manage inconsistencies
	Getters []Getter
	Setters []Setter
	Listers []Lister
}

func (c *Database) ListRPC(ctx context.Context, req *rpc.ListReq) (*rpc.ListResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var resp map[string]string
	for _, lister := range c.Listers {
		resp = lister.List(req.Key)
		if len(resp) > 0 {
			break
		}
	}
	response := &rpc.ListResp{
		Entries: resp,
		Code:    CoreCode_OKAY,
	}
	return response, nil
}

func (c *Database) SetRPC(ctx context.Context, req *rpc.SetReq) (*rpc.SetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	for _, setter := range c.Setters {
		setter.Set(req.Key, req.Val)
	}
	response := &rpc.SetResp{
		Code: CoreCode_OKAY,
	}
	return response, nil
}

func (c *Database) GetRPC(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	resp := ""
	for _, getter := range c.Getters {
		resp = getter.Get(req.Key)
		if resp != "" {
			break
		}
	}
	response := &rpc.GetResp{
		Val:  resp,
		Code: CoreCode_OKAY,
	}
	return response, nil
}

// NewDatabase will create a new database, and always initialize an in-memory cache.
func NewDatabase() *Database {
	memDB := mem.NewDatabase()
	return &Database{
		mtx: sync.Mutex{},
		Getters: []Getter{
			memDB,
		},
		Setters: []Setter{
			memDB,
		},
		Listers: []Lister{
			memDB,
		},
	}
}
