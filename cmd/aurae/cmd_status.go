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
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/crypto"
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
			
			foundErr := false
			con := printer.NewConsole("Status")

			kv := printer.NewKeyValueTable("")

			kv.AddKeyValue("Version", aurae.Version)

			// Private Key
			key, err := crypto.KeyFromPath(run.key)
			if err != nil {
				kv.AddKeyValue("Private Key", err)
			} else {
				kv.AddKeyValue("Private Key", run.key)
			}

			// Name (peer)
			self := peer.Self(key)
			kv.AddKeyValue("Name", self.Name.String())

			// Client
			auraeClient := client.NewClient()
			err = auraeClient.ConnectSocket(run.socket)
			if err != nil {
				foundErr = true
				kv.AddKeyValueErr("Client", err)
			} else {
				kv.AddKeyValue("Client", "Connected")
			}

			// Daemon
			_, err = auraeClient.Status(context.Background(), &rpc.StatusReq{})
			if err != nil {
				foundErr = true
				kv.AddKeyValueErr("Daemon", err)
			} else {
				kv.AddKeyValue("Daemon", "Running")
			}

			if foundErr {
				con.Title = "NOT READY"
			} else {
				con.Title = "READY"
			}
			con.AddPrinter(kv)
			con.PrintStdout()

			return nil
		},
	}
}
