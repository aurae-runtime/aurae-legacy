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
	"context"
	"fmt"
	crack "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/system"
	"github.com/sirupsen/logrus"
)

const (
	Path string = "/var/run/firecracker.socket"
	Name string = "firecracker"
)

var _ system.Socket = &Firecracker{}

type Firecracker struct {
	path   string
	name   string
	client *crack.Client
}

func (f *Firecracker) Adopt() error {
	client := crack.NewClient(f.path, logrus.NewEntry(logrus.New()), false)
	resp, err := client.GetMmds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to adopt firecracker: %v", err)
	}
	logrus.Infof("%v", resp.Payload)
	return nil
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

func NewFirecracker() system.Socket {
	return &Firecracker{
		path: Path,
		name: Name,
	}
}
