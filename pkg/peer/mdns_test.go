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

//func TestPeerMDNSLookup(t *testing.T) {
//
//	// + 0
//	p3 := NewPeer(name.New("beeps@nivenly.com"))
//	err := p3.Establish(context.Background(), 0)
//	if err != nil {
//		t.Errorf("unable to establish beeps@nivenly.com: %v", err)
//	}
//	err = p3.HandshakeServe()
//	if err != nil {
//		t.Errorf("unable to start handshake server on p3: %v", err)
//	}
//
//	// + 1
//	p4 := NewPeer(name.New("boops@nivenly.com"))
//	err = p4.Establish(context.Background(), 1)
//	if err != nil {
//		t.Errorf("unable to establish boops@nivenly.com: %v", err)
//	}
//	err = p4.Handshake(p3.Host().ID())
//	if err != nil {
//		t.Errorf("unable to handshake: %v", err)
//	}
//
//}
