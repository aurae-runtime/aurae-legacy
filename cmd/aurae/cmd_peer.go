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
	"github.com/kris-nova/aurae/pkg/crypto"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/printer"
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
		Flags: GlobalFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "servicename",
				Aliases:     []string{"svc"},
				Destination: &run.servicename,
			},
		}),
		Action: func(c *cli.Context) error {
			key, err := crypto.KeyFromPath(run.key)
			if err != nil {
				return err
			}
			svc := peer.NewPeer(name.New(run.servicename), key)
			host, err := svc.Establish()
			if err != nil {
				return err
			}

			con := printer.NewConsole("")
			tabPeer := printer.NewTable(fmt.Sprintf("PeerID: %s", host.ID()))
			prettyField := tabPeer.NewField("Peer")
			pubKeyTypeField := tabPeer.NewField("Public Key Type")
			pubKeyField := tabPeer.NewField("Public Key Type")
			for _, p := range host.Peerstore().PeersWithAddrs() {
				prettyField.AddValue(p.Pretty())
				pk, err := p.ExtractPublicKey()
				if err != nil {
					continue
				}
				pubKeyTypeField.AddValue(pk.Type())
				rawKeyData, err := pk.Raw()
				if err == nil {
					pubKeyField.AddValue(string(rawKeyData))
				}
			}
			tabPeer.AddField(prettyField)
			tabPeer.AddField(pubKeyField)
			tabPeer.AddField(pubKeyTypeField)
			con.AddTable(tabPeer)

			// Connect to the default service
			selfStream, err := host.NewStream(context.Background(), svc.ID())
			if err != nil {
				return err
			}
			tabStream := printer.NewTable(fmt.Sprintf("Stream: %s", selfStream.Protocol()))
			con.AddTable(tabStream)

			con.PrintStdout()
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
