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
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/core/memfs"
	"github.com/kris-nova/aurae/rpc"
	"sync"
)

var _ rpc.CoreServiceServer = &PathDatabase{}

const (
	CoreCode_OKAY  int32 = 0
	CoreCode_ERROR int32 = 1
	CoreCode_EMPTY int32 = 2
)

// PathDatabase is the core data store structure for Aurae.
//
// This structure couples the concepts of "path" filesystem to a
// key value store.
//
// If you store "boops" it becomes "/boops"
// If you store "beeps/boops" it becomes "/beeps/boops"
//
// The path paradigm allows for a less effecient but more queryable data set.
// This also allows us to expose the database as a mountable system over the socket.
type PathDatabase struct {
	mtx sync.Mutex
	rpc.UnimplementedCoreServiceServer
	// TODO we will need to prioritize our getters and setters
	// TODO we need to manage inconsistencies
	Getters []Getter
	Setters []Setter
	Listers []Lister
}

func (c *PathDatabase) ListRPC(ctx context.Context, req *rpc.ListReq) (*rpc.ListResp, error) {
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

// SetRPC is liable to mutate data!
func (c *PathDatabase) SetRPC(ctx context.Context, req *rpc.SetReq) (*rpc.SetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var key, val, mutPath string
	key = req.Key
	val = req.Val

	// Path mutation!
	mutPath = common.Path(key)

	// Ignore empty paths
	if mutPath == "/" {
		return &rpc.SetResp{
			Code: CoreCode_EMPTY,
		}, nil
	}

	// Path mutation!

	for _, setter := range c.Setters {
		setter.Set(mutPath, val)
	}
	response := &rpc.SetResp{
		Code: CoreCode_OKAY,
	}
	return response, nil
}

func (c *PathDatabase) GetRPC(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
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

// NewPathDatabase will create a new PathDatabase, and always initialize an in-memory cache.
func NewPathDatabase() *PathDatabase {
	memDB := memfs.NewDatabase()
	return &PathDatabase{
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
