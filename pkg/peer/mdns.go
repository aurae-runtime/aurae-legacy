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
	"github.com/sirupsen/logrus"
	"sync"
)

type NameService struct {
	AddrInfoCh chan peer.AddrInfo
	Records    map[name.Name]peer.AddrInfo
	mtx        sync.Mutex
}

func (s *NameService) HandlePeerFound(info peer.AddrInfo) {
	logrus.Debugf("Peer discovery: %v", info)
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.AddrInfoCh <- info
}

func NewNameService() *NameService {
	dns := &NameService{
		AddrInfoCh: make(chan peer.AddrInfo),
		Records:    make(map[name.Name]peer.AddrInfo),
		mtx:        sync.Mutex{},
	}
	return dns
}
