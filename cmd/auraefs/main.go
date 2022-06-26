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
	aurfs "github.com/kris-nova/aurae/aurafs"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/runtime"
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
			{
				Name:      "mount",
				Usage:     "Mount an auraeFS system.",
				UsageText: `auraefs mount <path>`,
				Flags: GlobalFlags([]cli.Flag{
					&cli.StringFlag{
						Name:        "socket",
						Aliases:     []string{"sock"},
						Destination: &run.socket,
						Value:       runtime.DefaultSocketLocationLinux,
					},
				}),
				Action: func(c *cli.Context) error {
					path := c.Args().Get(0)
					if path == "" {
						return fmt.Errorf("usage: auraefs mount <path>")
					}
					auraeClient := client.NewClient(run.socket)
					err := auraeClient.Connect()
					if err != nil {
						return err
					}
					fs := aurfs.NewAuraeFS(path, auraeClient)
					err = fs.Mount()
					if err != nil {
						return err
					}
					fs.Runtime()
					return nil
				},
			},
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
