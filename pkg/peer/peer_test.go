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
	"testing"
)

func TestPeerToPeerConnectSingle(t *testing.T) {
	a := NewPeer("a", nil)
	b := NewPeer("b", nil)
	var err error
	_, err = a.Connect()
	if err != nil {
		t.Errorf("unable to connect: %v", err)
	}
	_, err = b.Connect()
	if err != nil {
		t.Errorf("unable to connect: %v", err)
	}
	_, err = a.NewSafeConnection()
	if err != nil {
		t.Errorf("unable to net.conn: %v", err)
	}
	_, err = b.NewSafeConnection()
	if err != nil {
		t.Errorf("unable to net.conn: %v", err)
	}
}
