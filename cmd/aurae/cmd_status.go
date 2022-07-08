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
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/printer"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

const (
	StatusReady   string = "ready"
	StatusError   string = "error"
	StatusAlive   string = "alive"
	StatusUnknown string = "unknown"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show aurae status.",
		UsageText: `aurae status <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			input := c.Args().Get(0)
			var instance *peer.Peer
			isLocal := false
			instance = peer.Self()
			if input == "" {
				isLocal = true
			}

			err := instance.Establish(context.Background(), 0)
			if err != nil {
				return fmt.Errorf("unable to establish host: %s: %v", instance.Name.String(), err)
			}

			if isLocal {
				// Sock
			} else {
				// Peer
			}

			var daemonErr error

			statusResp, err := instance.Status(context.Background(), &rpc.StatusReq{})
			if err != nil {
				daemonErr = fmt.Errorf("unable to get status: %v", err)
			}
			//
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
			return con.PrintStdout()
		},
	}
}
