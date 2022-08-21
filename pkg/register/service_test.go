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
	"github.com/kris-nova/aurae/providers/tsocket"
	"github.com/kris-nova/aurae/rpc/rpc"
	"net"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove(tsocket.Path)
	conn, err := net.Listen("unix", tsocket.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	os.Exit(m.Run())
}

func TestService_AbandonSocket(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	// Assert the code
	if resp.Code != common.ResponseCode_OKAY {
		t.Errorf("failure testing tsocket: %v", resp.Message)
	}
	// Now close it
	resp2, err := svc.AbandonSocket(context.Background(), &rpc.AbandonSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	if resp2.Code != common.ResponseCode_OKAY {
		t.Errorf("failure testing tsocket: %v", resp2.Message)
	}

	// Now close it and confirm this is rejected because nothing exists
	resp3, err := svc.AbandonSocket(context.Background(), &rpc.AbandonSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	if resp3.Code != common.ResponseCode_REJECT {
		t.Errorf("failure testing tsocket: %v", resp2.Message)
	}
}

func TestService_AdoptSocket(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	// Assert the code
	if resp.Code != common.ResponseCode_OKAY {
		t.Errorf("failure testing tsocket: %v", resp.Message)
	}
}

func TestService_AdoptSocketSadPath(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                "/bad/socket/path",
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	// Assert the code
	if resp.Code != common.ResponseCode_ERROR {
		t.Errorf("expected error from service")
	}
}

func TestService_AdoptSocketSadName(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: "bad/name",
	})
	if err != nil {
		t.Errorf("error running service: %v", err)
	}
	// Assert the code
	if resp.Code != common.ResponseCode_ERROR {
		t.Errorf("expected error from service")
	}
}
