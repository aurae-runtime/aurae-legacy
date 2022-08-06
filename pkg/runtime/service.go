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
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
)

var _ rpc.RuntimeServer = &Service{}

type Service struct {
	rpc.UnimplementedRuntimeServer
}

func (s Service) Run(ctx context.Context, in *rpc.RunReq) (*rpc.RunResp, error) {
	// Spoof a response for now
	logrus.Debugf("Running image: %s", in.Name)
	return &rpc.RunResp{
		Code:    1,
		Message: "Success",
	}, nil
}

func (s Service) Status(ctx context.Context, in *rpc.StatusReq) (*rpc.StatusResp, error) {
	// Spoof a response for now
	return &rpc.StatusResp{
		Code:    1,
		Message: "Success",
		ProcessTable: map[string]*rpc.Process{
			"/init": {},
		},
		Field: "RUNNING",
	}, nil
}

func NewService() *Service {
	return &Service{}
}
