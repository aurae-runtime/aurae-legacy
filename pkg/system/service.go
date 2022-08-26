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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/system"
)

var _ aurae.SystemServer = &Service{}

type Service struct {
	aurae.UnimplementedSystemServer
}

func (s *Service) Status(ctx context.Context, in *aurae.StatusRequest) (*aurae.StatusResponse, error) {
	raw, err := system.AuraeInstance().Encapsulate()
	if err != nil {
		return &aurae.StatusResponse{
			Code:    common.ResponseCode_OKAY,
			Message: err.Error(),
		}, nil
	}
	return &aurae.StatusResponse{
		AuraeInstanceEncapsulated: string(raw),
		Code:                      common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
