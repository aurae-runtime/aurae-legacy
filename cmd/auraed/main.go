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
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/crypto"
	"github.com/kris-nova/aurae/pkg/daemon"
	"os"
	"time"

	"github.com/kris-nova/aurae"
	"github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

var run = &RuntimeOptions{}

type RuntimeOptions struct {
	verbose    bool
	socket     string
	localStore string
	key        string
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
		Usage:     "Aurae runtime daemon.",
		UsageText: `auraed <options>`,
		Commands:  []*cli.Command{
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
		},
		Flags:                GlobalFlags([]cli.Flag{}),
		EnableBashCompletion: true,
		HideHelp:             false,
		HideVersion:          false,

		Action: func(c *cli.Context) error {
			// Arbitrary (non-error) pre load
			Preloader()
			d := daemon.New(run.socket, run.localStore)
			return d.Run()
		},
	}

	var err error

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
		&cli.StringFlag{
			Name:        "key",
			Aliases:     []string{"s"},
			Destination: &run.key,
			Value:       fmt.Sprintf("%s/.ssh/%s", common.HomeDir(), crypto.DefaultAuraePrivateKeyName),
		},
		&cli.StringFlag{
			Name:        "socket",
			Aliases:     []string{"sock"},
			Destination: &run.socket,
			Value:       daemon.DefaultSocketLocationLinux,
			EnvVars: []string{
				"AURAE_SOCKET",
			},
		},
		&cli.StringFlag{
			Name:        "local",
			Aliases:     []string{"store"},
			Destination: &run.localStore,
			Value:       daemon.DefaultLocalStateLocationLinux,
		},
	}

	for _, gf := range g {
		c = append(c, gf)
	}
	return c
}

// Preloader will run for ALL commands, and is used
// to system the daemon environments of the program.
func Preloader() {
	/* Flag parsing */
	if run.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	run.socket = common.Expand(run.socket)
	run.localStore = common.Expand(run.localStore)
}
