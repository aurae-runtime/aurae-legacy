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
)

var _ aurae.LocalRuntimeServer = &Service{}

type Service struct {
	aurae.LocalRuntimeServer
}

func (s *Service) ReadStderr(ctx context.Context, in *aurae.ReadStderrRequest) (*aurae.ReadStderrResponse, error) {
	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &aurae.ReadStderrResponse{
			Message: "CapRunProcess not registered. Missing capability.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	// Polymorphism here
	if runtimeExecutor, ok := executor.(aurae.LocalRuntimeServer); ok {
		return runtimeExecutor.ReadStderr(ctx, in) // Run the implementation
	}

	return &aurae.ReadStderrResponse{
		Message: fmt.Sprintf("Unable to cast exector to RunProcessCapability"),
		Code:    common.ResponseCode_ERROR,
	}, nil
}

func (s *Service) ReadStdout(ctx context.Context, in *aurae.ReadStdoutRequest) (*aurae.ReadStdoutResponse, error) {
	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &aurae.ReadStdoutResponse{
			Message: "CapRunProcess not registered. Missing capability.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	// Polymorphism here
	if runtimeExecutor, ok := executor.(aurae.LocalRuntimeServer); ok {
		return runtimeExecutor.ReadStdout(ctx, in) // Run the implementation
	}

	return &aurae.ReadStdoutResponse{
		Message: fmt.Sprintf("Unable to cast exector to RunProcessCapability"),
		Code:    common.ResponseCode_ERROR,
	}, nil
}

func (s *Service) RunProcess(ctx context.Context, in *aurae.RunProcessRequest) (*aurae.RunProcessResponse, error) {
	executor := system.AuraeInstance().CapRunProcess
	if executor == nil {
		return &aurae.RunProcessResponse{
			Message: "CapRunProcess not registered. Missing capability.",
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

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
