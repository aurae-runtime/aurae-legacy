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
	"context"
	"github.com/kris-nova/aurae/pkg/name"
	"testing"
)

func TestPeerToPeerHandshake(t *testing.T) {

	ctx := context.Background()
	var err error

	p1 := NewPeer(name.New("p1"))
	err = p1.Establish(ctx, 0)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	p2 := NewPeer(name.New("p1"))
	err = p2.Establish(ctx, 1)
	if err != nil {
		t.Errorf("unable to establish p2: %v", err)
	}

	err = p1.HandshakeServe()
	if err != nil {
		t.Errorf("unable to start handshake serve on p1: %v", err)
	}
	err = p2.Handshake(p1.Host().ID())
	if err != nil {
		t.Errorf("unable to handshake from p2 -> p1: %v", err)
	}

}

func TestPeerToPeerHandshakeMultiplex(t *testing.T) {

	ctx := context.Background()
	var err error

	p1 := NewPeer(name.New("p1"))
	err = p1.Establish(ctx, 0)
	if err != nil {
		t.Errorf("unable to establish p1: %v", err)
	}
	p2 := NewPeer(name.New("p1"))
	err = p2.Establish(ctx, 1)
	if err != nil {
		t.Errorf("unable to establish p2: %v", err)
	}

	// p2 -> p1
	err = p1.HandshakeServe()
	if err != nil {
		t.Errorf("unable to start handshake serve on p1: %v", err)
	}
	err = p2.Handshake(p1.Host().ID())
	if err != nil {
		t.Errorf("unable to handshake from p2 -> p1: %v", err)
	}

	// p1 -> p2
	err = p2.HandshakeServe()
	if err != nil {
		t.Errorf("unable to start handshake serve on p2: %v", err)
	}
	err = p1.Handshake(p2.Host().ID())
	if err != nil {
		t.Errorf("unable to handshake from p1 -> p2: %v", err)
	}

}
