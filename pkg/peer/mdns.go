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
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/libp2p/go-libp2p-core/peer"
)

// TODO We need to understand if we want to host this in a (the?) DHT

type NameService struct {
	Peers map[string]*Record
}

type Record struct {
	ID   peer.ID
	Name name.Name
	Info peer.AddrInfo
}

// TODO We need to figure out where to host state

func (s *NameService) HandlePeerFound(info peer.AddrInfo) {
	//logrus.Infof("Peer discovery: %v", info)
	//logrus.Warnf("Stateless name mapping with raw addr info: %s", info.String())
	name := name.New(info.String())
	//logrus.Infof("Peer registered in DNS registry: [%s]", name.String())
	s.Peers[name.String()] = &Record{
		Name: name,
		Info: info,
		ID:   info.ID,
	}
}

func NewNameService() *NameService {
	return &NameService{
		Peers: make(map[string]*Record),
	}
}

// GetAddrInfo is a knock off of glibc getaddrinfo.
//
// Hopefully without the bugs and controversy.
func (s *NameService) GetAddrInfo(name string) *Record {
	if record, ok := s.Peers[name]; ok {
		return record
	}
	return nil
}
