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

package crypto

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"io/ioutil"
)

const (
	DefaultPrivateKeyName      string = "id_ed25519"
	DefaultPublicKeyName       string = "id_ed25519.pub"
	DefaultAuraePrivateKeyName string = "id_aurae"
	DefaultAuraePublicKeyName  string = "id_aurae.pub"
)

func KeyFromPath(path string) (crypto.PrivKey, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// TODO We should check the keys and support all the SSH keys
	return crypto.UnmarshalEd25519PrivateKey(bytes)
	//return crypto.UnmarshalPrivateKey(bytes)
}
