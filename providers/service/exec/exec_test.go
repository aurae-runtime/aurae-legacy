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

package exec

import (
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"testing"
)

func TestExec_RunProcess_ls_stdout(t *testing.T) {

	// Run
	r := &aurae.RunProcessRequest{
		ExecutableCommand: "ls",
	}

	var pid int32

	e := NewExec()
	if l, ok := e.(aurae.LocalRuntimeServer); !ok {
		t.Errorf("exec does not implement LocalRuntimeServer")
	} else {
		resp, err := l.RunProcess(nil, r)
		if err != nil {
			t.Errorf("unable to run: %v", err)
		}
		if resp.Code != common.ResponseCode_OKAY {
			t.Errorf("unable to run: %v", err)
		}
		pid = resp.PID
		if pid == 0 {
			t.Errorf("unable to find PID from process")
			t.FailNow()
		}
		procResp, err := l.GetProcessMeta(nil, &aurae.GetProcessMetaRequest{
			PID: pid,
		})
		if err != nil {
			t.Errorf("unable to get process meta: %v", err)
			t.FailNow()
		}
		if procResp.Code == common.ResponseCode_ERROR {
			t.Errorf("unable to get process meta, error code")
			t.FailNow()
		}
		if procResp.Stdout != "" {
			t.Logf("stdout: \n")
			t.Logf("%v", procResp.Stdout)
		}
		if procResp.Stdout != "" {
			t.Logf("stderr: \n")
			t.Logf("%v", procResp.Stderr)
		}

	}

}

func TestExec_RunProcess_ls(t *testing.T) {

	// Run
	r := &aurae.RunProcessRequest{
		ExecutableCommand: "ls",
	}

	e := NewExec()
	if l, ok := e.(aurae.LocalRuntimeServer); !ok {
		t.Errorf("exec does not implement LocalRuntimeServer")
	} else {
		resp, err := l.RunProcess(nil, r)
		if err != nil {
			t.Errorf("unable to run: %v", err)
		}
		if resp.Code != common.ResponseCode_OKAY {
			t.Errorf("unable to run: %v", err)
		}
	}
}
