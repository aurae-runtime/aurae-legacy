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
	"context"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/pkg/daemon"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	MountPath string = "/tmp/aurae"
)

func TestMain(m *testing.M) {

	// Bootstrap testing environment
	go func() {
		os.MkdirAll("/tmp/aurae", 0755)
	}()

	// Start a daemon

	d := daemon.New("/run/aurae.test.sock", "/tmp/aurae.test")
	ch := make(chan bool)
	go func() {
		ch <- true
		err := d.Run(context.TODO())
		if err != nil {
			logrus.Errorf("Error running daemon: %v", err)
			os.Exit(1)
		}
	}()

	// Create a client
	c := client.NewClient()
	err := c.ConnectSocket("/run/aurae.test.sock")
	if err != nil {
		logrus.Errorf("Error establishing client: %v", err)
		os.Exit(1)
	}

	// Mount the filesystem
	fs := NewAuraeFS("/tmp/aurae", c)
	go func() {
		err := fs.Mount()
		if err != nil {
			logrus.Errorf("Error mounting filesystem: %v", err)
			os.Exit(1)
		}
		ch <- true
		fs.Runtime()
	}()

	<-ch
	<-ch
	time.Sleep(1 * time.Second)
	exitCode := m.Run()

	// Clean up testing environment

	os.Exit(exitCode)
}

func TestInitialFilesystem(t *testing.T) {
	// cd /tmp/aurae
	err := ioutil.WriteFile(filepath.Join(MountPath, "testFile"), []byte("testData"), 0644)
	if err != nil {
		t.Errorf("unable to write testFile to auraefs: %v", err)
	}
	data, err := ioutil.ReadFile(filepath.Join(MountPath, "testFile"))
	if err != nil {
		t.Errorf("unable to read testFile to auraefs: %v", err)
	}
	if string(data) != "testData" {
		t.Errorf("Data IO failure. Expected: %s, Actual: %s", "testData", string(data))
	}
}
