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
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/name"
	peer "github.com/kris-nova/aurae/pkg/peer"
	peer2 "github.com/libp2p/go-libp2p-core/peer"
	"github.com/urfave/cli/v2"
	"os"
)

// Peer is the command for p2p
func Peer() *cli.Command {
	return &cli.Command{
		Name:      "peer",
		Usage:     "Create a new peer address and do nothing. Useful for debugging peer to peer.",
		UsageText: `aurae peer <options>`,
		Flags: GlobalFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "name",
				Aliases:     []string{"n"},
				Destination: &run.name,
				Value:       common.Self,
			},
		}),
		Action: func(c *cli.Context) error {
			//key, err := crypto.KeyFromPath(run.key)
			//if err != nil {
			//	return err
			//}
			var nameStr string
			if run.name == "" {
				hostname, err := os.Hostname()
				if err != nil {
					return fmt.Errorf("unable to calculate hostname: %v", err)
				}
				nameStr = hostname
			} else {
				nameStr = run.name
			}
			p := peer.NewPeer(name.New(nameStr))
			err := p.Establish(context.Background(), 1)
			if err != nil {
				return err
			}
			p.HandshakeServe()
			select {}
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
					hostname, err := os.Hostname()
					if err != nil {
						return fmt.Errorf("unable to calculate hostname: %v", err)
					}
					p := peer.NewPeer(name.New(hostname))
					err = p.Establish(context.Background(), 0)
					if err != nil {
						return err
					}
					peerID, _ := peer2.Decode(input)
					return p.Handshake(peerID)
				},
			},
		},
	}
}
