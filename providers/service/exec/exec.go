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
	"context"
	"fmt"
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/system"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	Name               string = "exec"
	CacheLengthSeconds int    = 1024
)

var _ system.Service = &Exec{}
var _ aurae.LocalRuntimeServer = &Exec{}

type Exec struct {
	aurae.LocalRuntimeServer
	sync.Mutex
	name  string
	cache map[int32]*ProcessMeta
}

type ProcessMeta struct {
	pid        int32
	x          *exec.Cmd
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser
}

func (e *Exec) Name() string {
	return e.name
}

func (e *Exec) Status() *system.ServiceStatus {
	return &system.ServiceStatus{}
}

func (e *Exec) Start() error {
	if system.AuraeInstance().CapRunProcess == nil {
		system.AuraeInstance().CapRunProcess = e
	} else {
		return fmt.Errorf("CapRunProcess already registered")
	}
	return nil
}

func (e *Exec) Stop() error {
	return nil
}

func NewExec() system.Service {
	return &Exec{
		cache: make(map[int32]*ProcessMeta),
		name:  Name,
		Mutex: sync.Mutex{},
	}
}

func (e *Exec) ReadStdout(ctx context.Context, in *aurae.ReadStdoutRequest) (*aurae.ReadStdoutResponse, error) {
	length := in.Length
	pid := in.PID
	e.Lock()
	defer e.Unlock()
	if procMeta, ok := e.cache[pid]; !ok {
		return &aurae.ReadStdoutResponse{
			Code:    common.ResponseCode_REJECT,
			Message: fmt.Sprintf("Pid %d not found in table.", pid),
		}, nil
	} else {
		pipe := procMeta.stdoutPipe
		buf := make([]byte, length)
		n, err := pipe.Read(buf)
		if err != nil {
			return &aurae.ReadStdoutResponse{
				Code:    common.ResponseCode_ERROR,
				Message: fmt.Sprintf("Unable to read bytes from pipe: %v", err),
			}, nil
		}
		// Success
		return &aurae.ReadStdoutResponse{
			PID:     pid,
			Size:    int32(n),
			Data:    string(buf),
			Code:    common.ResponseCode_OKAY,
			Message: common.ResponseMsg_Success,
		}, nil
	}
	return nil, nil
}

func (e *Exec) ReadStderr(ctx context.Context, in *aurae.ReadStderrRequest) (*aurae.ReadStderrResponse, error) {
	length := in.Length
	pid := in.PID
	e.Lock()
	defer e.Unlock()
	if procMeta, ok := e.cache[pid]; !ok {
		return &aurae.ReadStderrResponse{
			Code:    common.ResponseCode_REJECT,
			Message: fmt.Sprintf("Pid %d not found in table.", pid),
		}, nil
	} else {
		pipe := procMeta.stderrPipe
		buf := make([]byte, length)
		n, err := pipe.Read(buf)
		if err != nil {
			return &aurae.ReadStderrResponse{
				Code:    common.ResponseCode_ERROR,
				Message: fmt.Sprintf("Unable to read bytes from pipe: %v", err),
			}, nil
		}
		// Success
		return &aurae.ReadStderrResponse{
			PID:     pid,
			Size:    int32(n),
			Data:    string(buf),
			Code:    common.ResponseCode_OKAY,
			Message: common.ResponseMsg_Success,
		}, nil
	}
	return nil, nil
}

func (e *Exec) RunProcess(ctx context.Context, in *aurae.RunProcessRequest) (*aurae.RunProcessResponse, error) {
	execStart := in.ExecutableCommand // Nod to systemd :)
	spl := strings.Split(execStart, " ")

	// Break args apart
	var first string
	var args []string
	if len(spl) > 1 {
		first, args = spl[0], spl[1:]
	} else {
		first = execStart
	}

	// TODO We need to spend a lot of time evaluating this specific point in the code
	// TODO We need to spend a lot of time evaluating stdout, stderr buffers
	// TODO We need to spend a lot of time evaluating execve() and various arch/os types
	eCmd := exec.Command(first, args...)

	// Pipes
	stdoutPipe, err := eCmd.StdoutPipe()
	if err != nil {
		return &aurae.RunProcessResponse{
			Code: common.ResponseCode_ERROR,
		}, fmt.Errorf("unable to stdoutPipe: %v", err)
	}
	stderrPipe, err := eCmd.StderrPipe()
	if err != nil {
		return &aurae.RunProcessResponse{
			Code: common.ResponseCode_ERROR,
		}, fmt.Errorf("unable to stderrPipe: %v", err)
	}

	err = eCmd.Start()
	if err != nil {
		return &aurae.RunProcessResponse{
			Code:    common.ResponseCode_ERROR,
			Message: fmt.Sprintf("unable to start process: %v", err),
		}, fmt.Errorf("unable to start process: %v", err)
	}

	// Cache our process in the "process table"
	e.pidCache(eCmd, stdoutPipe, stderrPipe)

	return &aurae.RunProcessResponse{
		Code:    common.ResponseCode_OKAY,
		Message: "Started.",
		PID:     int32(eCmd.Process.Pid),
	}, nil

}

func (e *Exec) pidCache(x *exec.Cmd, stdoutPipe, stderrPipe io.ReadCloser) {
	pid := x.Process.Pid
	e.Lock()
	e.cache[int32(pid)] = &ProcessMeta{
		x:          x,
		stdoutPipe: stdoutPipe,
		stderrPipe: stderrPipe,
	}
	e.Unlock()

	go func() {
		time.Sleep(time.Duration(CacheLengthSeconds) * time.Second)
		x.Wait() // Wait will close the FDs "for our convenience", so wait and then call this!
		delete(e.cache, int32(pid))

	}()
}
