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

var runRun = &RuntimeRunOptions{}

type RuntimeRunOptions struct {
	command string
}

func Run() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run aurae processes, containers, and VMs.",
		UsageText: `aurae run`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Subcommands: []*cli.Command{
			{
				Name:      "process",
				Usage:     "Run aurae process.",
				UsageText: `aurae run process <arguments>`,
				Flags: GlobalFlags([]cli.Flag{
					&cli.StringFlag{
						Name:        "command",
						Aliases:     []string{"c", "x"},
						Value:       "",
						Destination: &runRun.command,
					},
				}),
				Action: func(c *cli.Context) error {
					Preloader()
					cmd := runRun.command
					if cmd == "" {
						return fmt.Errorf("usage: aurae run process -c 'command' <arguments>")
					}

					//spl := strings.Split(cmd, " ")
					//if len(spl) < 1 {
					//	return fmt.Errorf("usage: aurae run process -c 'command' <arguments>")
					//}
					//first, args := spl[0], spl[1:]

					ctx := context.Background()
					auraeClient := client.NewClient()
					err := auraeClient.ConnectSocket(run.socket)
					if err != nil {
						return err
					}
					resp, err := auraeClient.RunProcess(ctx, &aurae.RunProcessRequest{
						ExecutableCommand: cmd,
						Description:       "", //TODO Plumb description flag
						User:              "", //TODO Lookup user/plumb user string
					})
					if err != nil {
						return err
					}
					pid := resp.PID
					printer.PrintStdout("Run Process", resp)
					stdoutResp, err := auraeClient.ReadStdout(ctx, &aurae.ReadStdoutRequest{
						PID:    pid,
						Length: 1024 * 1024,
					})
					if err != nil {
						return err
					}
					if stdoutResp.Size != 0 {
						fmt.Println(stdoutResp.Data)
						stdoutResp.Data = "" // Unset the data
					}
					printer.PrintStdout("Stdout", stdoutResp)

					stderrResp, err := auraeClient.ReadStderr(ctx, &aurae.ReadStderrRequest{
						PID:    pid,
						Length: 1024 * 1024,
					})
					if err != nil {
						return err
					}
					if stderrResp.Size != 0 {
						fmt.Println(stderrResp.Data)
						stderrResp.Data = "" // Unset the data
					}
					printer.PrintStdout("Stderr", stderrResp)
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
