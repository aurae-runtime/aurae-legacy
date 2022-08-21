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

func TestService_AdoptSocket(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: tsocket.Name,
	})
	if err != nil {
		t.Errorf("unable to adopt socket: %v", err)
	}
	t.Logf("Adopted tsocket: %v", resp.Message)
}

func TestService_AdoptSocketSadPath(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                "/bad/socket/path",
		UniqueComponentName: tsocket.Name,
	})
	if err == nil {
		t.Errorf("expected error adopting bad socket")
	}
	t.Logf("Rejected bad socket: %v", resp.Message)
}

func TestService_AdoptSocketSadName(t *testing.T) {
	svc := NewService()
	resp, err := svc.AdoptSocket(context.Background(), &rpc.AdoptSocketRequest{
		Path:                tsocket.Path,
		UniqueComponentName: "bad/name",
	})
	if err == nil {
		t.Errorf("expected error adopting bad socket")
	}
	t.Logf("Rejected bad socket: %v", resp.Message)
}
