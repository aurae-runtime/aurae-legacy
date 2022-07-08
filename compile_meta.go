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

package aurae

import "fmt"

const (
	unknown string = "UNKNOWN"
)

var (
	Name        string = unknown
	Version     string = unknown
	Copyright   string = unknown
	License     string = unknown
	AuthorName  string = unknown
	AuthorEmail string = unknown
)

func Banner() string {
	var banner string
	banner += fmt.Sprintf("\n")
	banner += fmt.Sprintf(" ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓\n")
	banner += fmt.Sprintf(" ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗  ┃\n")
	banner += fmt.Sprintf(" ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗ ┃\n")
	banner += fmt.Sprintf(" ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║ ┃\n")
	banner += fmt.Sprintf(" ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║ ┃\n")
	banner += fmt.Sprintf(" ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║ ┃\n")
	banner += fmt.Sprintf(" ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝ ┃\n")
	banner += fmt.Sprintf(" ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
	banner += fmt.Sprintf(" Created by: %s <%s>\n", AuthorName, AuthorEmail)
	return banner
}
