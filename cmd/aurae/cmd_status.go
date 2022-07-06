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
	"github.com/kris-nova/aurae/pkg/printer"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

const (
	StatusReady   string = "ready"
	StatusError   string = "error"
	StatusAlive   string = "alive"
	StatusUnknwon string = "unknown"
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

			// Variables
			var auraeClient *client.Client
			auraeClient = client.NewClient()
			var daemonErr error

			// Initalize Peers and Clients
			//var p *peer.Peer
			//if input == "" || input == common.Self {
			//	p = peer.Self()
			//	err := auraeClient.ConnectSocket(run.socket)
			//	if err != nil {
			//		daemonErr = fmt.Errorf("unable to dial self socket: %v", err)
			//	}
			//} else {
			//	p = peer.NewPeer(name.New(input))
			//	err := p.Establish(context.Background(), 0)
			//	if err != nil {
			//		daemonErr = fmt.Errorf("unable to establish peer: %v", err)
			//	}
			//	err = auraeClient.ConnectPeer(p, "") // TODO dial to peer
			//	if err != nil {
			//		return fmt.Errorf("unable to dial peer: %v", err)
			//	}
			//}

			err := auraeClient.ConnectSocket(run.socket)
			if err != nil {
				daemonErr = fmt.Errorf("unable to dial self socket: %v", err)
			}

			// Get Status
			if auraeClient.RuntimeClient == nil {
				return fmt.Errorf("unable to join aurae net: %v", daemonErr)
			}

			statusResp, err := auraeClient.Status(context.Background(), &rpc.StatusReq{})
			if err != nil {
				daemonErr = fmt.Errorf("unable to get status: %v", err)
			}

			// Table 1
			con := printer.NewConsole("Status: Alive")
			t1 := printer.NewKeyValueTable("")

			if daemonErr != nil {
				t1.AddKeyValue("Daemon", StatusError)
				t1.AddKeyValueErr("Daemon Error", daemonErr)
				con.Title = "Status: Error"
			} else {
				t1.AddKeyValue("Daemon", StatusReady)
			}

			con.AddPrinter(t1)

			// Process Table
			if statusResp != nil && len(statusResp.ProcessTable) > 1 {
				t2 := printer.NewTable("Processes")
				nameField := t2.NewField("Name")
				procField := t2.NewField("Process")
				for name, proc := range statusResp.ProcessTable {
					nameField.AddValue(name)
					procField.AddValue(proc.Status)
				}
				con.AddPrinter(t2)
			}
			con.PrintStdout()
			return nil
		},
	}
}
