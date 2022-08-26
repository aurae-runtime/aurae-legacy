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
	"github.com/kris-nova/aurae/providers/service/exec"
	"github.com/kris-nova/aurae/system"
	"testing"
)

func TestService_RunProcess_Exec(t *testing.T) {

	// Establish Aurae with Exec for process management
	system.AuraeInstance().CapRunProcess = exec.NewExec()

	// New Service
	svc := NewService()
	resp, err := svc.RunProcess(context.Background(), &aurae.RunProcessRequest{
		ExecutableCommand: "ls /",
	})
	if err != nil {
		t.Errorf("unable to run process: %v", err)
		t.FailNow()
	}
	if resp.Code == common.ResponseCode_ERROR {
		t.Errorf("unable to run process: %v", resp.Message)
		t.FailNow()
	}

	t.Logf("Success: Pid: %d", resp.PID)

}
