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
	"github.com/kris-nova/aurae/pkg/runtime"
	"github.com/kris-nova/aurae/rpc"
	"os"
	"time"

	"github.com/kris-nova/aurae"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var run = &RuntimeOptions{}

type RuntimeOptions struct {
	verbose bool
	socket  string
}

func main() {
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "Show the version.",
	}
	app := &cli.App{
		Name:     aurae.Name,
		Version:  aurae.Version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  aurae.AuthorName,
				Email: aurae.AuthorEmail,
			},
		},
		Copyright: aurae.Copyright,
		HelpName:  aurae.Copyright,
		Usage:     "Simple, secure distributed system for application teams.",
		UsageText: `aurae <options> <command>`,
		Commands: []*cli.Command{
			// ----------------------------------------
			// [ STUB ]
			//{
			//	Name:      "NAME",
			//	Aliases:   []string{"NA", "N"},
			//	Usage:     "WHAT IT DOES",
			//	UsageText: `aurae <options> <command>`,
			//	Flags:     GlobalFlags([]cli.Flag{}),
			//	Action: func(c *cli.Context) error {
			//		return nil
			//	},
			//},
			// ----------------------------------------
			//{
			//	Name:      "status",
			//	Usage:     "View aurae status.",
			//	UsageText: `aurae <options> status`,
			//	Flags:     GlobalFlags([]cli.Flag{}),
			//	Action: func(c *cli.Context) error {
			//		printers.SimpleStdout(status.Now())
			//		return nil
			//	},
			//},
			{
				Name:      "get",
				Usage:     "Get aurae values.",
				UsageText: `aurae get <key>`,
				Flags: GlobalFlags([]cli.Flag{
					&cli.StringFlag{
						Name:        "socket",
						Aliases:     []string{"sock"},
						Destination: &run.socket,
						Value:       runtime.DefaultSocketLocationLinux,
					},
				}),
				Action: func(c *cli.Context) error {
					key := c.Args().Get(0)
					if key == "" {
						return fmt.Errorf("usage: aurae get <key>")
					}

					auraeClient := client.NewClient(run.socket)
					err := auraeClient.Connect()
					if err != nil {
						return err
					}
					getResp, err := auraeClient.GetRPC(context.Background(), &rpc.GetReq{
						Key: key,
					})
					if err != nil {
						return err
					}
					fmt.Println(getResp.Val)
					return nil
				},
			},
			{
				Name:      "set",
				Usage:     "Set aurae values.",
				UsageText: `aurae set <key> <value>`,
				Flags: GlobalFlags([]cli.Flag{
					&cli.StringFlag{
						Name:        "socket",
						Aliases:     []string{"sock"},
						Destination: &run.socket,
						Value:       runtime.DefaultSocketLocationLinux,
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

					auraeClient := client.NewClient(run.socket)
					err := auraeClient.Connect()
					if err != nil {
						return err
					}
					_, err = auraeClient.SetRPC(context.Background(), &rpc.SetReq{
						Key: key,
						Val: val,
					})
					if err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:      "list",
				Usage:     "List aurae values.",
				UsageText: `aurae list <key>`,
				Flags: GlobalFlags([]cli.Flag{
					&cli.StringFlag{
						Name:        "socket",
						Aliases:     []string{"sock"},
						Destination: &run.socket,
						Value:       runtime.DefaultSocketLocationLinux,
					},
				}),
				Action: func(c *cli.Context) error {
					key := c.Args().Get(0)
					if key == "" {
						return fmt.Errorf("usage: aurae get <key>")
					}

					auraeClient := client.NewClient(run.socket)
					err := auraeClient.Connect()
					if err != nil {
						return err
					}
					listResp, err := auraeClient.ListRPC(context.Background(), &rpc.ListReq{
						Key: key,
					})
					if err != nil {
						return err
					}
					for k, v := range listResp.Entries {
						fmt.Println(k, v)
					}
					return nil
				},
			},
			//{
			//	Name:      "run",
			//	Usage:     "Run an Application.",
			//	UsageText: `aurae <options> runtime`,
			//	Flags:     GlobalFlags([]cli.Flag{}),
			//	Action: func(c *cli.Context) error {
			//		image := c.Args().Get(0)
			//		if image == "" {
			//			return fmt.Errorf("usage: aurae run <image>")
			//		}
			//		conn, err := grpc.Dial(fmt.Sprintf("passthrough:///unix://%s", "/run/aurae.sock"), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
			//		if err != nil {
			//			return err
			//		}
			//		client := rpc.NewAuraeFSClient(conn)
			//
			//		nodeResp, err := client.GetRPC(context.Background(), &rpc.GetReq{
			//			Key: shape.FileEtcNameNode,
			//		})
			//		if err != nil {
			//			return err
			//		}
			//		domainResp, err := client.GetRPC(context.Background(), &rpc.GetReq{
			//			Key: shape.FileEtcNameDomain,
			//		})
			//		if err != nil {
			//			return err
			//		}
			//
			//		app := api.NewApplication(image, nodeResp.Val, domainResp.Val)
			//
			//		logrus.Infof("Saving new app: %s", app.Name)
			//
			//		//return fs.Set()
			//
			//		_, err = client.SetRPC(context.Background(), &rpc.SetReq{
			//			Key: filepath.Join(shape.DirApp, app.Name),
			//			Val: object.PrettyPrint(app),
			//		})
			//
			//		return err
			//	},
			//},
			//{
			//	Name:      "daemon",
			//	Usage:     "Run aurae as a runtime daemon.",
			//	UsageText: `aurae <options> runtime`,
			//	Flags: GlobalFlags([]cli.Flag{
			//		&cli.StringFlag{
			//			Name:        "mountpoint",
			//			Aliases:     []string{"m", "mount"},
			//			Destination: &run.mountpoint,
			//			Value:       "/aurae",
			//		},
			//		&cli.StringFlag{
			//			Name:        "socket",
			//			Aliases:     []string{"s", "sock"},
			//			Destination: &run.socket,
			//			Value:       "/run/aurae.sock",
			//		},
			//	}),
			//	Action: func(c *cli.Context) error {
			//		daemon := runtime.New(run.mountpoint, run.socket)
			//		return daemon.Run()
			//	},
			//},
			//{
			//	Name:      "peer",
			//	Usage:     "Peer this aurae instance to another.",
			//	UsageText: ` <options> status`,
			//	Flags:     GlobalFlags([]cli.Flag{}),
			//	Action: func(c *cli.Context) error {
			//		// TODO Peer
			//		printers.SimpleStdout(status.Now())
			//		return nil
			//	},
			//},
			//{
			//	Name:      "init",
			//	Usage:     "Initialize this aurae instance.",
			//	UsageText: `aurae <options> init <domain>`,
			//	Flags: GlobalFlags([]cli.Flag{
			//		&cli.StringFlag{
			//			Name:        "node",
			//			Aliases:     []string{"n", "nodename"},
			//			Destination: &run.nodename,
			//			Value:       "",
			//		},
			//	}),
			//	Action: func(c *cli.Context) error {
			//		domain := c.Args().Get(0)
			//		if domain == "" {
			//			return fmt.Errorf("usage: aurae init <domain>")
			//		}
			//		if run.nodename == "" {
			//			nodename, err := os.Hostname()
			//			if err != nil {
			//				return fmt.Errorf("unable to find nodename from hostname: %v", err)
			//			}
			//			run.nodename = nodename
			//		}
			//		// Initialize the local system
			//		fs, err := local.NewFromPath(run.root)
			//		if err != nil {
			//			return err
			//		}
			//		err = fs.Initialize(&shape.Config{
			//			DomainName: domain,
			//			NodeName:   run.nodename,
			//		})
			//		if err != nil {
			//			return err
			//		}
			//		printers.SimpleStdout(status.Now())
			//		return nil
			//	},
			//},
		},
		Flags:                GlobalFlags([]cli.Flag{}),
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,

		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
	}

	var err error

	// Load environment variables
	err = Environment()
	if err != nil {
		logrus.Error(err)
		os.Exit(99)
	}

	// Arbitrary (non-error) pre load
	Preloader()

	// Runtime
	err = app.Run(os.Args)
	if err != nil {
		logrus.Errorf(err.Error())
		os.Exit(-1)
	}
}

func GlobalFlags(c []cli.Flag) []cli.Flag {
	g := []cli.Flag{
		&cli.BoolFlag{
			Name:        "verbose",
			Aliases:     []string{"v"},
			Destination: &run.verbose,
		},
	}
	for _, gf := range g {
		c = append(c, gf)
	}
	return c
}

// Preloader will run for ALL commands, and is used
// to system the runtime environments of the program.
func Preloader() {
	/* Flag parsing */
	if run.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
