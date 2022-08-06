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
	"github.com/urfave/cli/v2"
)

func Capability() *cli.Command {
	return &cli.Command{
		Name: "capability",
		Aliases: []string{
			"cap",
		},
		Usage:     "The Aurare Capabilities API",
		UsageText: `aurae run <query>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			cli.ShowSubcommandHelpAndExit(c, 0)
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name: "get",
				Aliases: []string{
					"cap",
				},
				Usage:     "The Aurare Capabilities API",
				UsageText: `aurae cap get <options>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					cli.ShowSubcommandHelpAndExit(c, 0)
					return nil
				},
			},
			{
				Name: "set",
				Aliases: []string{
					"cap",
				},
				Usage:     "The Aurare Capabilities API",
				UsageText: `aurae cap set <options>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					cli.ShowSubcommandHelpAndExit(c, 0)
					return nil
				},
			},
		},
	}
}
