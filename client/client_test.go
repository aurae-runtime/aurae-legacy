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

package client

import (
	"context"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"testing"
)

func TestPeer2PeerConnect(t *testing.T) {
	var err error
	p1 := peer.Self()
	err = p1.Establish(context.Background(), 0)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	c1 := NewClient()

	p2 := peer.NewPeer(name.New("p2"))
	err = p2.Establish(context.Background(), 1)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	err = c1.ConnectPeer(p1, p2.Host().ID())
	if err != nil {
		t.Errorf("unable to connect client p1 -> p2: %v", err)
	}
}
