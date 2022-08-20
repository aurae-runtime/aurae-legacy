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
			"system",
		},
		Usage:     "The Aurare Capabilities API",
		UsageText: `aurae run <query>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			Preloader()
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name: "get",
				Aliases: []string{
					"system",
				},
				Usage:     "The Aurare Capabilities API",
				UsageText: `aurae system get <options>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					cli.ShowSubcommandHelpAndExit(c, 0)
					return nil
				},
			},
			{
				Name: "set",
				Aliases: []string{
					"system",
				},
				Usage:     "The Aurare Capabilities API",
				UsageText: `aurae system set <options>`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					cli.ShowSubcommandHelpAndExit(c, 0)
					return nil
				},
			},
		},
	}
}
