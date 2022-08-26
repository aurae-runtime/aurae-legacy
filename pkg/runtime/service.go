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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/system"
	"github.com/sirupsen/logrus"
)

var _ aurae.RuntimeServer = &Service{}

type Service struct {
	aurae.UnimplementedRuntimeServer
}

func (s *Service) RunProcess(ctx context.Context, in *aurae.RunProcessRequest) (*aurae.RunProcessResponse, error) {

	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &aurae.RunProcessResponse{
			Message: "CapRunProcess not registered. Missing capability.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	logrus.Infof("Command:   %s", in.ExecutablePath)
	logrus.Infof("Arguments: %s", in.ExecutableArgs)

	// TODO Left off here
	// TODO We now need to be able to build polymorphic capabilities
	// TODO If exec provides CapRunProcess we should be able to cast it to a new type

	return &aurae.RunProcessResponse{
		Code: common.ResponseCode_ERROR,
	}, nil
}

func (s *Service) RunContainer(ctx context.Context, in *aurae.RunContainerRequest) (*aurae.RunContainerResponse, error) {
	return &aurae.RunContainerResponse{
		Code: common.ResponseCode_OKAY,
	}, nil
}

func (s *Service) RunVirtualMachine(ctx context.Context, in *aurae.RunVirtualMachineRequest) (*aurae.RunVirtualMachineResponse, error) {
	return &aurae.RunVirtualMachineResponse{
		Code: common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
