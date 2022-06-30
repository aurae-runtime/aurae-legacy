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
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

func Peer() *cli.Command {
	return &cli.Command{
		Name:      "peer",
		Usage:     "Peer aurae with other aurae instances.",
		UsageText: `aurae peer <optional-token>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			token := c.Args().Get(0)
			rootClient := client.NewClient(run.socket)
			err := rootClient.Connect()
			if err != nil {
				return err
			}
			if token != "" {
				// Client mode
				proxyResp, err := rootClient.LocalProxy(context.Background(), &rpc.LocalProxyReq{
					Hostname: "TODO",
					Token:    token,
				})
				if err != nil {
					return err
				}
				fmt.Println(proxyResp.Hostname)
				fmt.Println(proxyResp.Socket)
				fmt.Println(proxyResp.Message)
				fmt.Println(proxyResp.Code)
			} else {
				// Server mode
				peerResp, err := rootClient.PeerRequest(context.Background(), &rpc.PeerRequestReq{})
				if err != nil {
					return err
				}

				fmt.Println("Token:")
				fmt.Println(peerResp.Token)
			}
			return nil
		},
	}
}
