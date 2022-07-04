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
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/crypto"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/printer"
	"github.com/urfave/cli/v2"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show aurae status.",
		UsageText: `aurae status <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {

			con := printer.NewConsole("Status")

			kv := printer.NewKeyValueTable("")

			kv.AddKeyValue("Version", aurae.Version)

			// Find key
			key, err := crypto.KeyFromPath(run.key)
			if err != nil {
				kv.AddKeyValue("Private Key", err)
			} else {
				kv.AddKeyValue("Private Key", run.key)
			}
			self := peer.Self(key)
			kv.AddKeyValue("Name", self.Name.String())

			auraeClient := client.NewClient()
			err = auraeClient.ConnectSocket(run.socket)
			if err != nil {
				kv.AddKeyValue("Client", err)
			} else {
				kv.AddKeyValue("Client", "Connected")
			}
			con.AddPrinter(kv)
			con.PrintStdout()

			return nil
		},
	}
}
