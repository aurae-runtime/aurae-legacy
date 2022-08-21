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

package tsocket

import (
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/system"
)

var _ system.Socket = &TSocket{}

const (
	Path string = "/tmp/auraetest.sock"
	Name string = "tsocket"
)

type TSocket struct {
	path string
	name string
}

func (f *TSocket) Path() string {
	return f.path
}

func (f *TSocket) Name() string {
	return f.name
}

func (f *TSocket) Status() *system.SocketStatus {
	return &system.SocketStatus{
		Message: aurae.Unknown,
	}
}

func NewTSocket(path, name string) system.Socket {
	return &TSocket{
		path: path,
		name: name,
	}
}
