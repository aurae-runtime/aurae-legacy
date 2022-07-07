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
	"fmt"
	"github.com/kris-nova/aurae/pkg/name"
	"os"
	"sync"
	"testing"
)

var (
	p1  *Peer
	p2  *Peer
	mtx *sync.Mutex
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	mtx = &sync.Mutex{}

	p1 = NewPeer(name.New("p1"))
	err = p1.Establish(ctx, 0)
	if err != nil {
		fmt.Sprintf("unable to establish p1: %v", err)
		os.Exit(1)
	}
	p2 = NewPeer(name.New("p2"))
	err = p2.Establish(ctx, 1)
	if err != nil {
		fmt.Sprintf("unable to establish p2: %v", err)
		os.Exit(2)
	}
	os.Exit(m.Run())
}

func TestPeerToPeerHandshakeMultiplex(t *testing.T) {
	mtx.Lock()
	defer mtx.Unlock()
	var err error

	// p2 -> p1
	err = p1.HandshakeServe()
	if err != nil {
		t.Errorf("unable to start handshake serve on p1: %v", err)
		t.FailNow()
	}
	err = p2.Handshake(p1.Host().ID())
	if err != nil {
		t.Errorf("unable to handshake from p2 -> p1: %v", err)
		t.FailNow()
	}

	// p1 -> p2
	err = p2.HandshakeServe()
	if err != nil {
		t.Errorf("unable to start handshake serve on p2: %v", err)
		t.FailNow()
	}
	err = p1.Handshake(p2.Host().ID())
	if err != nil {
		t.Errorf("unable to handshake from p1 -> p2: %v", err)
		t.FailNow()
	}
}
