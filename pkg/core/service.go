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
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/sirupsen/logrus"
	"sync"
)

var _ rpc.DatabaseServer = &Service{}

const (
	CoreCode_OKAY   int32 = 0
	CoreCode_ERROR  int32 = 1
	CoreCode_EMPTY  int32 = 2
	CoreCode_REJECT int32 = 3
)

var (
	// getFromMemory controls if we allow Aurae to get from memory ONLY.
	getFromMemory bool = true
)

type Service struct {
	mtx sync.Mutex
	rpc.UnimplementedDatabaseServer

	store CoreServicer
}

func (c *Service) List(ctx context.Context, req *rpc.ListReq) (*rpc.ListResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	path := common.Path(req.Key) // Data mutation!

	resp := make(map[string]bool)
	rpcResp := make(map[string]*rpc.Node)

	respMem := memfs.List(path)     // MemFS implementation
	respState := c.store.List(path) // LocalState implementation

	if getFromMemory {
		resp = respMem
	} else {
		resp = respState
	}

	for name, isFile := range resp {
		rpcResp[name] = &rpc.Node{
			Name: name,
			File: isFile,
		}
	}

	response := &rpc.ListResp{
		Entries: rpcResp,
		Code:    CoreCode_OKAY,
	}
	return response, nil
}

// Set is liable to mutate data!
func (c *Service) Set(ctx context.Context, req *rpc.SetReq) (*rpc.SetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	path := common.Path(req.Key) // Data mutation!

	// Ignore empty paths
	if path == "/" {
		return &rpc.SetResp{
			Code: CoreCode_EMPTY,
		}, nil
	}

	memfs.Set(path, req.Val)      // MemFS implementation
	go c.store.Set(path, req.Val) // LocalState implementation

	response := &rpc.SetResp{
		Code: CoreCode_OKAY,
	}
	return response, nil
}

// Remove is irreversible!
//
// To empty the database pass "/"
func (c *Service) Remove(ctx context.Context, req *rpc.RemoveReq) (*rpc.RemoveResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	path := common.Path(req.Key) // Data mutation!
	memfs.Remove(path)           // MemFS implementation
	go c.store.Remove(path)      // LocalState implementation

	response := &rpc.RemoveResp{
		Code: CoreCode_OKAY,
	}
	return response, nil
}

func (c *Service) Get(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var resp string

	path := common.Path(req.Key)   // Data mutation!
	respMem := memfs.Get(path)     // MemFS implementation
	respState := c.store.Get(path) // LocalState implementation

	if getFromMemory {
		resp = respMem
	} else {
		resp = respState
	}

	// Check for corruption.
	if respState != respMem {
		// TODO We need to talk through the implications of this and handle corruption.
		// logrus.Warnf("State corruption detected. Local != Memory.")
	}

	response := &rpc.GetResp{
		Val:  resp,
		Code: CoreCode_OKAY,
	}
	return response, nil
}

// NewService will create a new Service, and always initialize an in-memory cache.
//
// NewService depends on a CoreServicer which is a persistent datastore. The simplest
// datastore is merely local.State which is basic disk IO on a given path.
func NewService(store CoreServicer) *Service {
	return &Service{
		store: store,
		mtx:   sync.Mutex{},
	}
}

func (c *Service) SetGetFromMemory(x bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	getFromMemory = x
	logrus.Debugf("GetFromMemory: %t", x)
}
