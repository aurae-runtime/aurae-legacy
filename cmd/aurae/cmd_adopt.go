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
	"github.com/kris-nova/aurae/rpc/rpc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Adopt() *cli.Command {
	return &cli.Command{
		Name:      "adopt",
		Usage:     "Adopt aurae components. EG: Firecracker.",
		UsageText: `aurae adopt`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Subcommands: []*cli.Command{
			{
				Name:      "socket",
				Usage:     "Adopt socket components.",
				UsageText: `aurae adopt socket <path> <component-unique-name>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					Preloader()
					path := c.Args().Get(0)
					if path == "" {
						return fmt.Errorf("usage: aurae adopt socket <path> <component-unique-name>")
					}
					name := c.Args().Get(1)
					if name == "" {
						return fmt.Errorf("usage: aurae adopt socket <path> <component-unique-name>")
					}
					ctx := context.Background()
					auraeClient := client.NewClient()
					err := auraeClient.ConnectSocket(run.socket)
					if err != nil {
						return err
					}
					resp, err := auraeClient.AdoptSocket(ctx, &rpc.AdoptSocketRequest{
						Path:                path,
						UniqueComponentName: name,
					})
					if err != nil {
						return err
					}
					logrus.Info(resp.Code)
					logrus.Info(resp.Message)
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowSubcommandHelpAndExit(c, 0)
			return nil
		},
	}
}
