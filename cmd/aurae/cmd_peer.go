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
	"github.com/kris-nova/aurae/pkg/crypto"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/urfave/cli/v2"
)

// Peer
//
// TODO this needs to go over the aurae socket!
func Peer() *cli.Command {
	return &cli.Command{
		Name:      "peer",
		Usage:     "Work with Aurae peers in the mesh.",
		UsageText: `aurae peer <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			key, err := crypto.KeyFromPath(run.key)
			if err != nil {
				return err
			}
			self := peer.Self(key)
			host, err := self.Establish()
			if err != nil {
				return err
			}
			fmt.Printf("Peer ID: %s\n", host.ID())
			fmt.Printf("Peerstore:\n")
			for _, p := range host.Peerstore().PeersWithAddrs() {
				fmt.Printf(" - %s\n", p.String())
			}
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:      "",
				Usage:     "",
				UsageText: ``,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {

					return nil
				},
			},
		},
	}
}
