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
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Run() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run a container image.",
		UsageText: `aurae run <query>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			Preloader()
			input := c.Args().Get(0)
			if input == "" {
				return fmt.Errorf("Usage: aure run <image>. \nempty container image string")
			}
			query, err := client.Query(input)
			if err != nil {
				return fmt.Errorf("invalid query: %s", err)
			}
			x := query.Client
			name := query.Name
			ctx := query.Context
			err = x.ConnectSocket(run.socket)
			if err != nil {
				return fmt.Errorf("unable to connect: %s", err)
			}
			logrus.Debug("Running: %s", name.String())

			runResp, err := x.Run(ctx, &rpc.RunReq{})
			if err != nil {
				return fmt.Errorf("unable to run: %s", err)
			}
			logrus.Infof(runResp.String())
			return nil
		},
	}
}
