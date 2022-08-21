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

package runtime

import (
	"context"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/kris-nova/aurae/system"
)

var _ rpc.RuntimeServer = &Service{}

type Service struct {
	rpc.UnimplementedRuntimeServer
}

func (s *Service) RunProcess(ctx context.Context, in *rpc.RunProcessRequest) (*rpc.RunProcessResponse, error) {

	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &rpc.RunProcessResponse{
			Message: "CapRunProcess not provided.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	// TODO Left off here
	// TODO We now need to be able to build polymorphic capabilities
	// TODO If exec provides CapRunProcess we should be able to cast it to a new type

	return &rpc.RunProcessResponse{
		Code: common.ResponseCode_OKAY,
	}, nil
}

func (s *Service) RunContainer(ctx context.Context, in *rpc.RunContainerRequest) (*rpc.RunContainerResponse, error) {
	return &rpc.RunContainerResponse{
		Code: common.ResponseCode_OKAY,
	}, nil
}

func (s *Service) RunVirtualMachine(ctx context.Context, in *rpc.RunVirtualMachineRequest) (*rpc.RunVirtualMachineResponse, error) {
	return &rpc.RunVirtualMachineResponse{
		Code: common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
