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
	"context"
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/multiformats/go-multistream"
	"github.com/sirupsen/logrus"
	"io"
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

func (p *Peer) Handshake(id peer.ID) error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.host.SetStreamHandler(AuraeStreamProtocol(), doHandshake)

	s, err := p.host.NewStream(context.Background(), id, AuraeStreamProtocol())
	if err != nil {
		if err == multistream.ErrNotSupported {
			return fmt.Errorf("unable to create handshake stream, handshake server not discovered: enable %s on remote peer", AuraeStreamProtocol())
		}
		return fmt.Errorf("unable to create new stream: %v", err)
	}

	_, err = s.Write([]byte(AuraeProtocolHandshakeRequest))
	if err != nil {
		return fmt.Errorf("handshake failure write: %v", err)
	}
	response, err := io.ReadAll(s)
	if err != nil {
		return fmt.Errorf("handshake failure read: %v", err)
	}
	if string(response) != AuraeProtocolHandshakeResponse {
		return fmt.Errorf("handshake failure validate: %s", string(response))
	}
	logrus.Infof("Aurae handshake: Success.")
	return nil
}

func (p *Peer) HandshakeServe() error {
	if !p.established {
		return fmt.Errorf("unable to stream, first establish in the mesh")
	}
	p.host.SetStreamHandler(AuraeStreamProtocol(), doHandshake)
	return nil
}

func doHandshake(s network.Stream) {
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
