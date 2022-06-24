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

package shape

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kris-nova/aurae/cap"
)

// Model will model the filesystem, however will not hydrate the filesystem
// with any data. This is just the model of aurae.
type Model struct {
	root         string
	Files        map[string]os.FileMode
	Directories  map[string]os.FileMode
	Capabilities []*cap.Capability
}

const (
	DirApp   string = "app"
	DirUsr   string = "usr"
	DirSvc   string = "svc"
	DirCap   string = "cap"
	DirEtc   string = "etc"
	DirInfra string = "infra"

	FileEtcEnabled    string = "etc/enabled"
	FileEtcNameNode   string = "etc/name.node"
	FileEtcNameDomain string = "etc/name.domain"
	FileauraeSock     string = "aurae.sock"
	FileInfraNodeSelf string = "infra/self"
)

// NewModel will create a new shape model from a root.
// A *Model is a utility object and should only be used
// to model a filesystem. A *Model should NEVER be used
// for interacting with shape in production!
func NewModel(root string) *Model {
	if !strings.Contains(root, "aurae") {
		root = filepath.Join(root, "aurae")
	}
	return &Model{
		root: root,

		// Directories will be joined with root
		Directories: map[string]os.FileMode{
			DirApp:   DefaultMode,
			DirCap:   DefaultMode,
			DirUsr:   DefaultMode,
			DirSvc:   DefaultMode,
			DirEtc:   DefaultMode,
			DirInfra: DefaultMode,
		},

		Files: map[string]os.FileMode{
			FileEtcNameNode:   DefaultMode,
			FileEtcNameDomain: DefaultMode,
			FileEtcEnabled:    DefaultMode,
			//FileauraeSock:      SocketMode,
			FileInfraNodeSelf: DefaultMode,
		},
	}
}
