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
	"github.com/kris-nova/aurae/rpc"
)

var _ rpc.CoreServiceServer = &Database{}
var _ Getter = &Database{}
var _ Setter = &Database{}

type Database struct {
	rpc.UnimplementedCoreServiceServer
}

func (c Database) Get(key string) string {
	//TODO implement me
	panic("implement me")
}

func (c Database) Set(key, value string) {
	//TODO implement me
	panic("implement me")
}
func (c Database) SetRPC(ctx context.Context, req *rpc.SetReq) (*rpc.SetResp, error) {
	response := &rpc.SetResp{}
	return response, nil
}

func (c Database) GetRPC(ctx context.Context, req *rpc.GetReq) (*rpc.GetResp, error) {
	response := &rpc.GetResp{}
	return response, nil
}

func NewDatabase() *Database {
	return &Database{}
}