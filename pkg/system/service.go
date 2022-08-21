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

package system

import (
	"context"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/kris-nova/aurae/system"
)

var _ rpc.SystemServer = &Service{}

type Service struct {
	rpc.UnimplementedSystemServer
}

func (s *Service) Status(ctx context.Context, in *rpc.StatusRequest) (*rpc.StatusResponse, error) {
	raw, err := system.AuraeInstance().Encapsulate()
	if err != nil {
		return &rpc.StatusResponse{
			Code:    common.ResponseCode_OKAY,
			Message: err.Error(),
		}, nil
	}
	return &rpc.StatusResponse{
		AuraeInstanceEncapsulated: string(raw),
		Code:                      common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
