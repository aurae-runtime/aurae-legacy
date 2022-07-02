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

package hostname

import (
	"fmt"
	"strings"
)

// Hostname represents a 3 part hostname.
//
// beeps@boops@computer.com
type Hostname struct {
	Sub    string
	Host   string
	Domain string
}

func New(any string) *Hostname {
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
		Domain: domain,
		Host:   host,
		Sub:    sub,
	}
}

func (h *Hostname) String() string {
	return fmt.Sprintf("%s@%s@%s", h.Sub, h.Host, h.Domain)
}
