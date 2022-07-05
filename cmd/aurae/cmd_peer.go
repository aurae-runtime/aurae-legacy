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
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/sirupsen/logrus"
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
			svc := peer.Self()
			_, err := svc.Establish()
			if err != nil {
				return err
			}
			addr := svc.Address()
			logrus.Infof("Address: %s", addr)
			for {
			}
		},
		Subcommands: []*cli.Command{
			{
				Name:      "to",
				Usage:     "Establish a connection to another peer.",
				UsageText: `aurae peer to <addr>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					addr := c.Args().Get(0)
					if addr == "" {
						return fmt.Errorf("usage: aurae peer to <addr>")
					}
					//key, err := crypto.KeyFromPath(run.key)
					//if err != nil {
					//	return err
					//}
					svc := peer.Self()
					_, err := svc.Establish()
					if err != nil {
						return err
					}
					stream, err := svc.ToPeerAddr(addr)
					if err != nil {
						return err
					}
					logrus.Infof("Connected! Stream: %s", stream.ID())
					return nil
				},
			},
		},
	}
}
