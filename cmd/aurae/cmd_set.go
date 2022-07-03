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
	"github.com/kris-nova/aurae/pkg/daemon"
	"github.com/kris-nova/aurae/rpc"
	"github.com/urfave/cli/v2"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "Set aurae values.",
		UsageText: `aurae set <key> <value>`,
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
			val := c.Args().Get(1)
			if key == "" {
				return fmt.Errorf("usage: aurae set <key> <value>")
			}
			if val == "" {
				return fmt.Errorf("usage: aurae set <key> <value>")
			}

			auraeClient := client.NewClient()
			err := auraeClient.ConnectSocket(run.socket)
			if err != nil {
				return err
			}
			_, err = auraeClient.Set(context.Background(), &rpc.SetReq{
				Key: key,
				Val: val,
			})
			if err != nil {
				return err
			}
			return nil
		},
	}
}
