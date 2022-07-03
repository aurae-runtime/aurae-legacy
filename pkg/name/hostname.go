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

package name

import (
	"fmt"
	"strings"
)

// Hostname represents a 3 part name.
//
// beeps@boops@computer.com
type Hostname struct {
	sub    string
	host   string
	domain string
}

func NewHostname(any string) *Hostname {
	domain := ""
	host := ""
	sub := ""
	spl := strings.Split(any, "@")
	if len(spl) >= 3 {
		domain = spl[2]
		host = spl[1]
		sub = spl[0]
	} else if len(spl) == 2 {
		host = spl[0]
		domain = spl[1]
	} else if len(spl) == 1 {
		host = spl[0]
	} else {
		host = any
	}
	return &Hostname{
		domain: domain,
		host:   host,
		sub:    sub,
	}
}

// Host is the FQN of the host following the
// 3 @ format.
//
// beeps@alice@nivenly.com
// hack@alice@nivenly.com
func (h *Hostname) Host() string {
	return fmt.Sprintf("%s@%s@%s", h.sub, h.host, h.domain)
}
