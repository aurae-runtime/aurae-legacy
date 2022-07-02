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
// TODO this needs to go over the aurae socket!
func Peer() *cli.Command {
	return &cli.Command{
		Name:      "peer",
		Usage:     "Work with Aurae peers in the mesh.",
		UsageText: `aurae peer <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			key, err := peer.KeyFromPath(run.key)
			if err != nil {
				return err
			}
			self := peer.Self(key)
			host, err := self.Connect()
			if err != nil {
				return err
			}
			logrus.Infof("Runtime ID: %s", host.ID())
			for _, p := range host.Peerstore().PeersWithAddrs() {
				fmt.Println(p.String())
			}
			return nil
		},
	}
}
