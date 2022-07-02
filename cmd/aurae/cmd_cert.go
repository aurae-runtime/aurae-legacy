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
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

func Cert() *cli.Command {
	return &cli.Command{
		Name:      "cert",
		Usage:     "Work with local TLS certificates.",
		UsageText: `aurae cert`,
		Flags:     GlobalFlags([]cli.Flag{}),
		Action: func(c *cli.Context) error {
			cli.ShowSubcommandHelp(c)
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:      "keygen",
				Usage:     "Generate new keypair (id_ed25519, id_ed25519.pub)",
				UsageText: `aurae cert`,
				Flags:     GlobalFlags([]cli.Flag{}),
				Action: func(c *cli.Context) error {
					privKey, pubKey, err := crypto.GenerateKeyPair(crypto.Ed25519, 2048)
					if err != nil {
						return err
					}
					privKeyBytes, err := privKey.Raw()
					if err != nil {
						return err
					}
					err = ioutil.WriteFile("id_ed25519", privKeyBytes, 0644)
					if err != nil {
						return err
					}
					pubKeyBytes, err := pubKey.Raw()
					if err != nil {
						return err
					}
					err = ioutil.WriteFile("id_ed25519.pub", pubKeyBytes, 0644)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}
}
