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
	acrypto "github.com/kris-nova/aurae/pkg/crypto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sirupsen/logrus"
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
					privKeyPath := fmt.Sprintf("%s/.ssh/%s", common.HomeDir(), acrypto.DefaultAuraePrivateKeyName)
					err = ioutil.WriteFile(privKeyPath, privKeyBytes, 0644)
					if err != nil {
						return err
					}
					logrus.Infof("Wrote: %s", privKeyPath)
					pubKeyBytes, err := pubKey.Raw()
					if err != nil {
						return err
					}
					pubKeyPath := fmt.Sprintf("%s/.ssh/%s", common.HomeDir(), acrypto.DefaultAuraePublicKeyName)
					err = ioutil.WriteFile(pubKeyPath, pubKeyBytes, 0644)
					if err != nil {
						return err
					}
					logrus.Infof("Wrote %s", pubKeyPath)
					return nil
				},
			},
		},
	}
}
