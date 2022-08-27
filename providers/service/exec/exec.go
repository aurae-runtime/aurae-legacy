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
	"github.com/sirupsen/logrus"
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

//func (e *Exec) GetProcessMeta(ctx context.Context, in *aurae.GetProcessMetaRequest) (*aurae.GetProcessMetaResponse, error) {
//
//	pid := in.PID
//	e.Lock()
//	if x, ok := e.cache[pid]; !ok {
//		return &aurae.GetProcessMetaResponse{
//			Code: common.ResponseCode_ERROR,
//		}, nil
//	} else {
//
//		var stdoutStr, stderrStr string
//		if x.stderr != nil {
//			stderrStr = x.stderr.String()
//		}
//		if x.stdout != nil {
//			stdoutStr = x.stdout.String()
//		}
//
//		return &aurae.GetProcessMetaResponse{
//			Stderr: stderrStr,
//			Stdout: stdoutStr,
//			Code:   common.ResponseCode_OKAY,
//		}, nil
//	}
//
//	e.Unlock()
//	return &aurae.GetProcessMetaResponse{}, nil
//}

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
		err := x.Wait()
		if err != nil {
			logrus.Warnf("error waiting for process in table: %v", err)
		}
		time.Sleep(time.Duration(CacheLengthSeconds) * time.Second)
		delete(e.cache, int32(pid))
	}()
}
