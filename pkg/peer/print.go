package peer

import (
	"fmt"
	"github.com/kris-nova/aurae/pkg/printer"
)

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

func (p *Peer) Print() error {
	con := printer.NewConsole("Peer: %s", p.Name.String())
	kv1 := printer.NewKeyValueTable("")
	kv1.AddKeyValue("Peer ID", p.Host.ID())
	kv1.AddKeyValue("ServiceName", p.Name.Service())
	kv1.AddKeyValue("HostName", p.Name.Host())
	for i, addr := range p.Host.Addrs() {
		key := fmt.Sprintf("Addr%d", i)
		kv1.AddKeyValue(key, addr.String())
	}
	con.AddPrinter(kv1)
	con.PrintStdout()
	return nil
}
