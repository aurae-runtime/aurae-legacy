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

package common

import (
	"fmt"
	"strings"
)

// ParseRawName will parse a raw name and provide its fully qualified domain
// name in return.
func ParseRawName(raw, node, domain string) string {
	spl := strings.Split(raw, "@")
	if len(spl) == 3 {
		if spl[1] != node {
			return raw
		}
		if spl[2] != domain {
			return raw
		}
		return fmt.Sprintf("%s@%s@%s", spl[0], node, domain)
	}
	if len(spl) == 2 {
		if spl[1] == node {
			return fmt.Sprintf("%s@%s@%s", spl[0], node, domain)
		}
		if spl[1] == domain {
			return fmt.Sprintf("%s@%s@%s", spl[0], node, spl[1])
		}
	}
	if len(spl) == 1 {
		return fmt.Sprintf("%s@%s@%s", spl[0], node, domain)
	}
	return fmt.Sprintf("%s@%s@%s", raw, node, domain)
}
