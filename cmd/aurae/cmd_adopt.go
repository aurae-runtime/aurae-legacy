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
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/printer"
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
					resp, err := auraeClient.AdoptSocket(ctx, &aurae.AdoptSocketRequest{
						Path:                path,
						UniqueComponentName: name,
					})
					if err != nil {
						return err
					}
					printer.PrintStdout("Adopt Socket", resp)
					return nil
				},
			},
			{
				Name:      "service",
				Usage:     "Adopt service components.",
				UsageText: `aurae adopt service <name>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					Preloader()
					name := c.Args().Get(0)
					if name == "" {
						return fmt.Errorf("usage: aurae adopt service <name>")
					}
					ctx := context.Background()
					auraeClient := client.NewClient()
					err := auraeClient.ConnectSocket(run.socket)
					if err != nil {
						return err
					}
					resp, err := auraeClient.AdoptService(ctx, &aurae.AdoptServiceRequest{
						UniqueComponentName: name,
					})
					if err != nil {
						return err
					}
					printer.PrintStdout("Adopt Service", resp)
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
