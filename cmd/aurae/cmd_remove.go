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
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/pkg/daemon"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

func Remove() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove aurae values.",
		UsageText: `aurae remove <key>`,
		Flags: GlobalFlags([]cli.Flag{
			&cli.StringFlag{
				Name:        "socket",
				Aliases:     []string{"sock"},
				Destination: &run.socket,
				Value:       daemon.DefaultSocketLocationLinux,
			},
		}),
		Action: func(c *cli.Context) error {
			key := c.Args().Get(0)
			if key == "" {
				return fmt.Errorf("usage: aurae remove <key>")
			}

			auraeClient := client.NewClient(run.socket)
			err := auraeClient.Connect()
			if err != nil {
				return err
			}
			removeResp, err := auraeClient.Remove(context.Background(), &rpc.RemoveReq{
				Key: key,
			})
			if err != nil {
				return err
			}

			if removeResp.Code != core.CoreCode_OKAY {
				return fmt.Errorf("unable to remove")
			}
			return nil
		},
	}
}