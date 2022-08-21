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
	Unknown string = "Unknown"
)

var (
	Version     string = Unknown
	Copyright   string = Unknown
	License     string = Unknown
	AuthorName  string = Unknown
	AuthorEmail string = Unknown
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
