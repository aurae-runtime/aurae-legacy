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
				err = instance.ConnectSock(run.socket)
				if err != nil {
					return fmt.Errorf("unable to connect to local host socket %s: %s: %v", run.socket, instance.Name.String(), err)
				}
			} else {
				// Peer
				err = instance.ConnectPeerString(input)
				if err != nil {
					return fmt.Errorf("unable to connect to peer grpc net %s: %s: %v", input, instance.Name.String(), err)
				}
			}

			//
			//// Variables
			//var auraeClient *client.Client
			//auraeClient = client.NewClient()
			var daemonErr error
			//
			//// Initalize Peers and Clients
			//var p *peer.Peer
			//if input == "" || input == common.Self {
			//	p = peer.Self()
			//	err := auraeClient.ConnectSocket(run.socket)
			//	if err != nil {
			//		daemonErr = fmt.Errorf("unable to dial self socket: %v", err)
			//	}
			//} else {
			//	//p = peer.Self()
			//	//err := p.Establish(context.Background(), 0)
			//	//if err != nil {
			//	//	daemonErr = fmt.Errorf("unable to establish peer: %v", err)
			//	//}
			//	//to, err := peer2peer.Decode(input)
			//	//if err != nil {
			//	//	return fmt.Errorf("unable to get peer id: %v", err)
			//	//}
			//	//err = auraeClient.ConnectPeer(p, to) // TODO dial to peer
			//	//if err != nil {
			//	//	return fmt.Errorf("unable to dial peer: %v", err)
			//	//}
			//}
			//
			//err := inst.ConnectSocket(run.socket)
			//if err != nil {
			//	daemonErr = fmt.Errorf("unable to dial self socket: %v", err)
			//}
			////
			//// Get Status
			//if auraeClient.RuntimeClient == nil {
			//	return fmt.Errorf("unable to join aurae net: %v", daemonErr)
			//}
			//

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
