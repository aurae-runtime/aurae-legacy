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
	"fmt"
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/system"
	"github.com/sirupsen/logrus"
)

var _ aurae.LocalRuntimeServer = &Service{}

type Service struct {
	aurae.LocalRuntimeServer
}

// TODO This can (and should) be simpler.

//type RunProcessCapability interface {
//	RunProcess(ctx context.Context, in *aurae.RunProcessRequest) (*aurae.RunProcessResponse, error)
//}

func (s *Service) RunProcess(ctx context.Context, in *aurae.RunProcessRequest) (*aurae.RunProcessResponse, error) {

	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &aurae.RunProcessResponse{
			Message: "CapRunProcess not registered. Missing capability.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	logrus.Infof("ExecutableCommand: %s", in.ExecutableCommand)
	logrus.Infof("Description      : %s", in.Description)

	// Polymorphism here
	if runtimeExecutor, ok := executor.(aurae.LocalRuntimeServer); ok {
		return runtimeExecutor.RunProcess(ctx, in) // Run the implementation
	}

	return &aurae.RunProcessResponse{
		Message: fmt.Sprintf("Unable to cast exector to RunProcessCapability"),
		Code:    common.ResponseCode_ERROR,
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
