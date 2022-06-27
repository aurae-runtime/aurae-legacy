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
}

func (c *PathDatabase) ListRPC(ctx context.Context, req *rpc.ListReq) (*rpc.ListResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	path := common.Path(req.Key) // Data mutation!

	resp := make(map[string]*rpc.Node)

	ls := memfs.List(path) // MemFS implementation
	
	// Copy the memfs.Node -> rpc.Node // TODO should we simplify this type?
	for name, node := range ls {
		file := false
		if node != nil {
			file = node.File
		}
		resp[name] = &rpc.Node{
			Name: name,
			File: file,
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

	path := common.Path(req.Key) // Data mutation!

	// Ignore empty paths
	if path == "/" {
		return &rpc.SetResp{
			Code: CoreCode_EMPTY,
		}, nil
	}

	memfs.Set(path, req.Val) // MemFS implementation

	response := &rpc.SetResp{
		Code: CoreCode_OKAY,
	}
	return response, nil
}

func (c *PathDatabase) GetRPC(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	path := common.Path(req.Key) // Data mutation!

	resp := memfs.Get(path) // MemFS implementation

	response := &rpc.GetResp{
		Val:  resp,
		Code: CoreCode_OKAY,
	}
	return response, nil
}

// NewPathDatabase will create a new PathDatabase, and always initialize an in-memory cache.
func NewPathDatabase() *PathDatabase {
	return &PathDatabase{
		mtx: sync.Mutex{},
	}
}
