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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/registry"
	"github.com/kris-nova/aurae/system"
	"github.com/sirupsen/logrus"
	"os"
)

var _ aurae.RegisterServer = &Service{}

type Service struct {
	aurae.UnimplementedRegisterServer
}

func (s *Service) AdoptService(ctx context.Context, in *aurae.AdoptServiceRequest) (*aurae.AdoptServiceResponse, error) {

	name := in.UniqueComponentName

	// Check if registered
	if newFunc, ok := registry.ServiceRegistry[name]; ok {

		// Check if already loaded
		a := system.AuraeInstance()
		if _, ok := a.ServiceComponents[name]; ok {
			return &aurae.AdoptServiceResponse{
				Message: fmt.Sprintf("Already registered: %s", name),
				Code:    common.ResponseCode_REJECT,
			}, nil
		}

		// Load this service
		serviceInstance := newFunc()
		err := serviceInstance.Start()
		if err != nil {
			return &aurae.AdoptServiceResponse{
				Message: fmt.Sprintf("Unable to adopt service: %s: %v", name, err),
				Code:    common.ResponseCode_ERROR,
			}, nil
		}
		a.ServiceComponents[name] = serviceInstance
	} else {
		return &aurae.AdoptServiceResponse{
			Message: fmt.Sprintf("Service not found in registry: %s", name),
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	logrus.Infof("Success. Registered service: %s", name)
	return &aurae.AdoptServiceResponse{
		Message: fmt.Sprintf("Success. Registered socket: %s", name),
		Code:    common.ResponseCode_OKAY,
	}, nil
}

func (s *Service) AbandonService(ctx context.Context, in *aurae.AbandonServiceRequest) (*aurae.AbandonServiceResponse, error) {
	name := in.UniqueComponentName

	a := system.AuraeInstance()
	if _, ok := a.SocketComponents[name]; ok {
		delete(a.SocketComponents, name)
		logrus.Infof("Success. Abandoned service: %s", name)
		return &aurae.AbandonServiceResponse{
			Message: fmt.Sprintf("Stopped service: %s", name),
			Code:    common.ResponseCode_OKAY,
		}, nil
	}
	return &aurae.AbandonServiceResponse{
		Message: fmt.Sprintf("Service not found in registry: %s", name),
		Code:    common.ResponseCode_REJECT,
	}, nil
}

func (s *Service) AbandonSocket(ctx context.Context, in *aurae.AbandonSocketRequest) (*aurae.AbandonSocketResponse, error) {
	name := in.UniqueComponentName

	a := system.AuraeInstance()
	if _, ok := a.SocketComponents[name]; ok {
		delete(a.SocketComponents, name)
		logrus.Infof("Success. Abandoned socket: %s", name)
		return &aurae.AbandonSocketResponse{
			Message: fmt.Sprintf("Closed socket: %s", name),
			Code:    common.ResponseCode_OKAY,
		}, nil
	}
	return &aurae.AbandonSocketResponse{
		Message: fmt.Sprintf("Socket not found in registry: %s", name),
		Code:    common.ResponseCode_REJECT,
	}, nil
}

func (s *Service) AdoptSocket(ctx context.Context, in *aurae.AdoptSocketRequest) (*aurae.AdoptSocketResponse, error) {
	name := in.UniqueComponentName
	path := in.Path

	// Check if a valid socket exists at this location
	stat, err := os.Stat(path)
	if err != nil {
		return &aurae.AdoptSocketResponse{
			Message: fmt.Sprintf("Unable to stat file type: %s: %v", path, err),
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	// Srwxr-xr-x    Example socket with "S" and other permissions
	// S---------    Example os.ModeSocket with "S" only
	if stat.Mode()&os.ModeSocket == 0 {
		return &aurae.AdoptSocketResponse{
			Message: fmt.Sprintf("File not of type socket: %s: %v", path, err),
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	// Check if already loaded
	a := system.AuraeInstance()
	if _, ok := a.SocketComponents[name]; ok {
		return &aurae.AdoptSocketResponse{
			Message: fmt.Sprintf("Already registered: %s", name),
			Code:    common.ResponseCode_REJECT,
		}, nil
	}

	// Check if exists in registry, adopt if found
	if newFunc, ok := registry.SocketRegistry[name]; ok {
		socketInstance := newFunc()
		err := socketInstance.Adopt()
		if err != nil {
			return &aurae.AdoptSocketResponse{
				Message: fmt.Sprintf("Unable to adopt socket: %s: %v", name, err),
				Code:    common.ResponseCode_ERROR,
			}, nil
		}
		a.SocketComponents[name] = socketInstance
	} else {
		return &aurae.AdoptSocketResponse{
			Message: fmt.Sprintf("Socket not found in registry: %s", name),
			Code:    common.ResponseCode_ERROR,
		}, nil
	}

	logrus.Infof("Success. Registered socket: %s [%s]", name, path)
	return &aurae.AdoptSocketResponse{
		Message: fmt.Sprintf("Success. Registered socket: %s [%s]", name, path),
		Code:    common.ResponseCode_OKAY,
	}, nil
}

func NewService() *Service {
	return &Service{}
}
