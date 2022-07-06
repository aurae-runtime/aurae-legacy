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

package main

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/printer"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show aurae status.",
		UsageText: `aurae status <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			input := c.Args().Get(0)
			if input == "" {
				input = common.Self
			}

			var auraeClient *client.Client
			auraeClient = client.NewClient()

			var p *peer.Peer
			if input == "" || input == common.Self {
				p = peer.Self()
				err := auraeClient.ConnectSocket(run.socket)
				if err != nil {
					return fmt.Errorf("unable to dial self socket: %v", err)
				}
			} else {
				p = peer.NewPeer(name.New(input))
				err := p.Establish(context.Background(), 1)
				if err != nil {
					return fmt.Errorf("unable to establish peer: %v", err)
				}
				err = auraeClient.ConnectPeer(p)
				if err != nil {
					return fmt.Errorf("unable to dial self socket: %v", err)
				}
			}

			statusResp, err := auraeClient.Status(context.Background(), &rpc.StatusReq{})
			if err != nil {
				return fmt.Errorf("unable to get status: %v", err)
			}

			con := printer.NewConsole("Status")
			t1 := printer.NewKeyValueTable("")
			t1.AddKeyValue("Field", statusResp.Field)
			t1.AddKeyValue("Code", statusResp.Code)
			t1.AddKeyValue("Message", statusResp.Message)

			t2 := printer.NewTable("Processes")
			nameField := t2.NewField("Name")
			procField := t2.NewField("Process")
			for name, proc := range statusResp.ProcessTable {
				nameField.AddValue(name)
				procField.AddValue(proc.Status)
			}

			con.AddPrinter(t1)
			con.AddPrinter(t2)
			con.PrintStdout()
			return nil
		},
	}
}
