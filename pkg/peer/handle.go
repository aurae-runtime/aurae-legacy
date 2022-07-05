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

package peer

import (
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const (
	AuraeStream              string = "/aurae"    // The official stream endpoint for Aurae
	AuraeStreamVersionFormat string = "/aurae/%s" // Format with the package version
)

func AuraeStreamProtocol() protocol.ID {
	auraeStreamProtocol := fmt.Sprintf(AuraeStreamVersionFormat, aurae.Version)
	ids := protocol.ConvertFromStrings([]string{auraeStreamProtocol})
	if len(ids) != 1 {
		panic("unable to find aurae protocol!")
	}
	return ids[0]
}

func HandlePeerConnectStream(s network.Stream) {
	logrus.Infof("Received new stream. Protocol: %s", s.Protocol())
	data, err := ioutil.ReadAll(s)
	if err != nil {
		logrus.Warnf("Unable to read network stream: %v", err)
	}
	logrus.Infof("Raw stream:")
	logrus.Infof("%v", data)
}
