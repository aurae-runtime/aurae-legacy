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
	"fmt"
	zpeer "github.com/kris-nova/aurae/pkg/peer/peer"
	"github.com/urfave/cli/v2"
)

// Peer
//
//
func Peer() *cli.Command {
	return &cli.Command{
		Name:      "peer",
		Usage:     "Create a new peer address and do nothing. Useful for debugging peer to peer.",
		UsageText: `aurae peer <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			//key, err := crypto.KeyFromPath(run.key)
			//if err != nil {
			//	return err
			//}
			zpeer.RunServer()
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:      "to",
				Usage:     "Establish a connection to another peer.",
				UsageText: `aurae peer to <addr>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					input := c.Args().Get(0)
					if input == "" {
						return fmt.Errorf("usage: aurae peer to <addr>")
					}
					zpeer.RunClient(input)

					return nil
				},
			},
		},
	}
}
