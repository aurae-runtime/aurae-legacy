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
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func Status() *cli.Command {
	return &cli.Command{
		Name:      "status",
		Usage:     "Show aurae status.",
		UsageText: `aurae status <options>`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			query := c.Args().Get(0)
			first, _, _ := common.PathSplit(query)
			var hostname string
			if first == "" {
				// Local mode
				hostname = common.Localhost
			} else {
				// Peer mode
				hostname = first
			}

			logrus.Infof("Host: %s", hostname)
			var clientToUse *client.Client
			auraeClient := client.NewClient()
			err := auraeClient.ConnectSocket(run.socket)
			if err != nil {
				return err
			}

			//if hostname != common.Localhost {
			//	peerClient, err := rootClient.NewPeer(hostname)
			//	if err != nil {
			//		return err
			//	}
			//	err = peerClient.Connect()
			//	if err != nil {
			//		return err
			//	}
			//	clientToUse = peerClient
			//} else {
			//	clientToUse = rootClient
			//}
			statusResp, err := clientToUse.Status(context.Background(), &rpc.StatusReq{})
			if err != nil {
				return err
			}
			fmt.Println(statusResp.ProcessTable)
			fmt.Println(statusResp.Code)
			fmt.Println(statusResp.Field)
			fmt.Println(statusResp.Message)
			return nil
		},
	}
}
