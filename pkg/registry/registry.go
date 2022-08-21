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

package registry

import (
	"github.com/kris-nova/aurae/providers/service/exec"
	"github.com/kris-nova/aurae/providers/socket/firecracker"
	"github.com/kris-nova/aurae/providers/socket/tsocket"
	"github.com/kris-nova/aurae/system"
)

type NewSocket func() system.Socket
type NewService func() system.Service

var (
	ServiceRegistry = map[string]NewService{

		// exec is a simple executor for processes using fork()
		"exec": exec.NewExec,
	}

	SocketRegistry = map[string]NewSocket{

		// tsocket is a test socket provider, useful for unit tests
		"tsocket": tsocket.NewTSocket,

		// firecracker is a microvm management system from Amazon
		"firecracker": firecracker.NewFirecracker,
	}
)
