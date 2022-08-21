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
	"fmt"
	"github.com/kris-nova/aurae/system"
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
