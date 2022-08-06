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
	"github.com/fatih/color"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/daemon"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

func List() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List aurae values.",
		UsageText: `aurae list <key>`,
		Flags: GlobalFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "socket",
				Aliases:     []string{"sock"},
				Destination: &run.socket,
				Value:       daemon.DefaultSocketLocationLinux,
			},
		}),
		Action: func(c *cli.Context) error {
			Preloader()
			key := c.Args().Get(0)
			if key == "" {
				return fmt.Errorf("usage: aurae list <key>")
			}

			auraeClient := client.NewClient()
			err := auraeClient.ConnectSocket(run.socket)
			if err != nil {
				return err
			}
			listResp, err := auraeClient.List(context.Background(), &rpc.ListReq{
				Key: key,
			})
			if err != nil {
				return err
			}
			for k, v := range listResp.Entries {
				if v.File {
					color.Green(k)
				} else {
					color.Blue(k)
				}
			}
			return nil
		},
	}
}
