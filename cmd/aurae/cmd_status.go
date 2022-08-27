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
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/printer"
	"github.com/kris-nova/aurae/system"
	"github.com/urfave/cli/v2"
)

const (
	StatusReady   string = "ready"
	StatusError   string = "error"
	StatusAlive   string = "alive"
	StatusUnknown string = "unknown"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show aurae status.",
		UsageText: `aurae status <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			Preloader()
			ctx := context.Background()
			auraeClient := client.NewClient()
			err := auraeClient.ConnectSocket(run.socket)
			if err != nil {
				return err
			}
			status, err := auraeClient.Status(ctx, &aurae.StatusRequest{})
			if err != nil {
				return err
			}
			if status.Code == common.ResponseCode_OKAY {
				auraeInstance, err := system.StringToAuraeSafe(status.AuraeInstanceEncapsulated)
				if err != nil {
					printer.PrintStdout("status", status)
					return fmt.Errorf("unable to marshal aurae status from remote: %v", err)
				}
				printer.PrintStdout("status", auraeInstance)
				return nil
			}
			printer.PrintStdout("status", status)
			return fmt.Errorf("error getting status: %v", status.Message)
			return nil
		},
	}
}
