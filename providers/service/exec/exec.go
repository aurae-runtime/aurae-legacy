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
	"bytes"
	"context"
	"fmt"
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/system"
	"os/exec"
	"strings"
)

const (
	Name string = "exec"
)

var _ system.Service = &Exec{}

type Exec struct {
	name string
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
		name: Name,
	}
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

	var stderr, stdout bytes.Buffer
	eCmd.Stdout = &stdout
	eCmd.Stderr = &stderr
	err := eCmd.Start()
	if err != nil {
		return &aurae.RunProcessResponse{
			Code:    common.ResponseCode_ERROR,
			Message: fmt.Sprintf("unable to start process: %v", err),
		}, fmt.Errorf("unable to start process: %v", err)
	}

	return &aurae.RunProcessResponse{
		Code:    common.ResponseCode_OKAY,
		Message: "Started.",
		PID:     int32(eCmd.Process.Pid),
	}, nil

}
