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
	"github.com/kris-nova/aurae/pkg/common"
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
		Usage:     "Learn about a peer in the mesh. IDs or service names accepted.",
		UsageText: `aurae peer <peer> <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {

			peerQuery := c.Args().Get(0)
			if peerQuery == "" {
				peerQuery = common.Self // No query, use self
			}

			key, err := crypto.KeyFromPath(run.key)
			if err != nil {
				return err
			}

			svc := peer.NewPeerServicename(peerQuery, key)
			_, err = svc.Establish()
			if err != nil {
				return err
			}

			return svc.Print()
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
