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

package runtime

import (
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/aurafs"
	"github.com/kris-nova/aurae/pkg/posix"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
)

// Daemon is an aurae systemd style daemon.
//
// The daemon will securely mount an shape and
// expose the filesystem over mTLS gRPC for secure
// (multi-tenant) IPC over encrypted gRPC in the local
// userspace.
//
// While this methodology does come at a small performance
// hit, I believe this risk is justified by the secure multi tenancy
// feature.
//
type Daemon struct {
	runtime    bool
	mountpoint string
	socket     string
}

func New(mountpoint, socket string) *Daemon {
	return &Daemon{
		mountpoint: mountpoint,
		runtime:    true,
		socket:     socket,
	}
}

func (d *Daemon) Run() error {
	// Start the Signal Handler
	quitCh := posix.SignalHandler()
	go func() {
		d.runtime = <-quitCh
	}()
	var err error

	// Dispatch events from the filesystem
	//err := d.Dispatch()
	//if err != nil {
	//	logrus.Errorf("Dispatch failure: %v", err)
	//	logrus.Errorf("Shutting down.")
	//	d.runtime = false
	//}

	fs := aurafs.NewAuraeFS(d.mountpoint)
	err = fs.Mount()
	if err != nil {
		return fmt.Errorf("unable to mount AuraeFS: %v", err)
	}

	// Load gRPC server

	// Clean this up after we understand the relationship with auraeFS and gRPC with mTLS
	logrus.Infof("Aurae runtime daemon. Version: %s", aurae.Version)
	logrus.Infof("Initializing Unix Domain Socket...")
	conn, err := net.Listen("unix", d.socket)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			logrus.Warningf("Attempting to clean socket from failure!")
			os.Remove(d.socket)
			conn, err = net.Listen("unix", d.socket)
		}
	}
	if err != nil {
		return fmt.Errorf("unable to start daemon: socket error: %v", err)
	} else {
		logrus.Infof("Success. Socket acquired.")
	}

	server := grpc.NewServer()
	rpc.RegisterAuraeFSServer(server, fs)

	logrus.Infof("Starting gRPC Server")
	go server.Serve(conn)
	for d.runtime {
	}
	logrus.Infof("Gracefully shutting aurae service.")
	server.GracefulStop()
	server.Stop()
	logrus.Infof("Gentle cleanup.")
	logrus.Infof("Safely exiting.")
	return nil
}
