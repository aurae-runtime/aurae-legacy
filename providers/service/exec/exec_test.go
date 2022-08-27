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

func TestExec_RunProcess_ls_stderr(t *testing.T) {
	r := &aurae.RunProcessRequest{
		ExecutableCommand: "mv /SHOULDNOTEXISTSERIOUSLYTHISSHOULDNOTEXIST /tmp/SOMETHING",
	}
	var pid int32
	e := NewExec()
	if l, ok := e.(aurae.LocalRuntimeServer); !ok {
		t.Errorf("exec does not implement LocalRuntimeServer")
		t.FailNow()
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

		readResp, err := l.ReadStderr(nil, &aurae.ReadStderrRequest{
			PID:    pid,
			Length: 1024,
		})
		if err != nil {
			t.Errorf("unable to read stdout pipe: %v", err)
			t.FailNow()
		}
		if readResp.Code != common.ResponseCode_OKAY {
			t.Errorf("failed response readstderr")
			t.FailNow()
		}
		t.Logf("%v", readResp.Data)
	}
}

func TestExec_RunProcess_ls_stdout_singlebyte(t *testing.T) {
	r := &aurae.RunProcessRequest{
		ExecutableCommand: "ls -la /tmp",
	}
	var pid int32
	e := NewExec()
	if l, ok := e.(aurae.LocalRuntimeServer); !ok {
		t.Errorf("exec does not implement LocalRuntimeServer")
		t.FailNow()
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

		readMax := 1024
		read := 0
		for {
			readResp, err := l.ReadStdout(nil, &aurae.ReadStdoutRequest{
				PID:    pid,
				Length: 1,
			})
			if err != nil {
				t.Errorf("unable to read stdout pipe: %v", err)
				t.FailNow()
			}
			if readResp.Code != common.ResponseCode_OKAY {
				t.Errorf("failed response readstdout: %s", readResp.Message)
				t.FailNow()
			}
			//t.Logf("%v", readResp.Data)
			read = read + int(readResp.Size)
			if read >= readMax {
				break
			}
		}
	}
}

func TestExec_RunProcess_ls_stdout(t *testing.T) {
	r := &aurae.RunProcessRequest{
		ExecutableCommand: "ls -la /tmp",
	}
	var pid int32
	e := NewExec()
	if l, ok := e.(aurae.LocalRuntimeServer); !ok {
		t.Errorf("exec does not implement LocalRuntimeServer")
		t.FailNow()
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

		readResp, err := l.ReadStdout(nil, &aurae.ReadStdoutRequest{
			PID:    pid,
			Length: 1024,
		})
		if err != nil {
			t.Errorf("unable to read stdout pipe: %v", err)
			t.FailNow()
		}
		if readResp.Code != common.ResponseCode_OKAY {
			t.Errorf("failed response readstdout")
			t.FailNow()
		}
		t.Logf("%v", readResp.Data)
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
