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
			p := zpeer.NewPeer()
			err := p.Establish(context.Background(), 1)
			if err != nil {
				return err
			}
			return p.Stream()
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
					p := zpeer.NewPeer()
					err := p.Establish(context.Background(), 0)
					if err != nil {
						return err
					}
					return p.To(input)
				},
			},
		},
	}
}
