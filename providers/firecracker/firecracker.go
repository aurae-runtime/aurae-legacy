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

package firecracker

import (
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/system"
)

const (
	Path string = "/var/run/firecracker.socket"
	Name string = "firecracker"
)

var _ system.Socket = &Firecracker{}

type Firecracker struct {
	path string
	name string
}

func (f *Firecracker) Close() error {
	return nil
}

func (f *Firecracker) Path() string {
	return f.path
}

func (f *Firecracker) Name() string {
	return f.name
}

func (f *Firecracker) Status() *system.SocketStatus {
	return &system.SocketStatus{
		Message: aurae.Unknown,
	}
}

func NewFirecracker(path, name string) system.Socket {
	return &Firecracker{
		path: path,
		name: name,
	}
}
