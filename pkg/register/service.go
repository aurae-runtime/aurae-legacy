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

package register

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/registry"
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/kris-nova/aurae/system"
)

var _ rpc.RegisterServer = &Service{}

type Service struct {
	rpc.UnimplementedRegisterServer
}

func (s *Service) AdoptSocket(ctx context.Context, in *rpc.AdoptSocketRequest) (*rpc.AdoptSocketResponse, error) {
	name := in.UniqueComponentName
	path := in.Path

	// Check if exists in registry
	if newFunc, ok := registry.SocketRegistry[name]; ok {
		a := system.AuraeInstance()
		a.Sockets[name] = newFunc(name, path)
	} else {
		return &rpc.AdoptSocketResponse{
			Message: fmt.Sprintf("Socket not found in registry: %s", name),
			Code:    common.ResponseCode_ERROR,
		}, nil
	}
	return &rpc.AdoptSocketResponse{
		Message: fmt.Sprintf("Success. Registered socket: %s [%s]", name, path),
		Code:    common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
