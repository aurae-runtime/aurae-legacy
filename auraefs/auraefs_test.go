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

package auraefs

import (
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/runtime"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// Bootstrap testing environment

	// Start a daemon
	d := runtime.New("/run/aurae.test.sock")
	go func() {
		err := d.Run()
		if err != nil {
			logrus.Errorf("Error running daemon: %v", err)
			os.Exit(1)
		}
	}()

	// Create a client
	c := client.NewClient("/run/aurae.test.sock")

	// Mount the filesystem
	fs := NewAuraeFS("/tmp/aurae", c)
	go func() {
		err := fs.Mount()
		if err != nil {
			logrus.Errorf("Error mounting filesystem: %v", err)
			os.Exit(1)
		}
		fs.Runtime()
	}()

	exitCode := m.Run()

	// Clean up testing environment

	os.Exit(exitCode)
}

func TestInitialFilesystem(t *testing.T) {
	// cd /tmp/aurae
	// echo "data" > /tmp/aurae/file
	// cat /tmp/aurae/file
	// verify "data"
	// rm file
	// ls /tmp/aurae
}
