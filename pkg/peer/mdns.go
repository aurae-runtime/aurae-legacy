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
	"github.com/libp2p/go-libp2p-core/network"
	"net"
	"time"
)

var _ net.Conn = &Conn{}

type Conn struct {
	std net.Conn
	p2p network.Conn
}

func NewConn() *Conn {
	return &Conn{
		std: &net.TCPConn{},
	}
}

func (c Conn) Read(b []byte) (n int, err error) {
	return c.std.Read(b)
}

func (c Conn) Write(b []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (c Conn) Close() error {
	//TODO implement me
	panic("implement me")
}

func (c Conn) LocalAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c Conn) RemoteAddr() net.Addr {
	//TODO implement me
	panic("implement me")
}

func (c Conn) SetDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c Conn) SetReadDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (c Conn) SetWriteDeadline(t time.Time) error {
	//TODO implement me
	panic("implement me")
}
