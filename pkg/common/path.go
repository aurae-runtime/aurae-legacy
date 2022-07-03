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
	"path"
	"strings"
)

// Path is a heavily tested (and heavily used) common Path sanitation function.
// This function is deterministic and will attempt to turn a key into a meaningful
// filesystem path.
func Path(key string) string {
	isMkdir := false
	if strings.HasSuffix(key, "/") {
		isMkdir = true
	}
	var ret string
	ret = strings.ReplaceAll(key, "\\", "")
	rawPieces := strings.Split(ret, "/")
	var cleanPieces []string
	for _, p := range rawPieces {
		p = strings.TrimSpace(p)
		cleanPieces = append(cleanPieces, p)
	}
	ret = strings.Join(cleanPieces, "/")
	ret = path.Clean(ret)
	ret = path.Join("/", ret)
	if isMkdir && ret != "/" {
		ret = fmt.Sprintf("%s/", ret)
	}
	return ret
}
