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
	"bufio"
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/sirupsen/logrus"
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

const (
	AuraeProtocolHandshakeError    string = "<--**<<ERROR>>**-->\n"
	AuraeProtocolHandshakeRequest  string = "<--**<<REQUEST>>**-->\n"
	AuraeProtocolHandshakeResponse string = "<--**<<RESPONSE>>**-->\n"
)

func Handshake(s network.Stream) {
	defer s.Close()
	buf := bufio.NewReader(s)
	handshakeStr, err := buf.ReadString('\n')
	if err != nil {
		logrus.Warnf("Handshake failure: %v", err)
		s.Write([]byte(AuraeProtocolHandshakeError)) // Error
		return
	}
	if handshakeStr != AuraeProtocolHandshakeRequest {
		s.Write([]byte(AuraeProtocolHandshakeError)) // Error
		return
	}
	s.Write([]byte(AuraeProtocolHandshakeResponse)) // Okay
	return
}
