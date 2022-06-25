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
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/pkg/posix"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
)

const (
	DefaultSocketLocationLinux string = "/run/aurae.sock"
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
	runtime bool
	socket  string
}

func New(socket string) *Daemon {
	return &Daemon{
		runtime: true,
		socket:  socket,
	}
}

func (d *Daemon) Run() error {

	// Step 1. Establish context in the logs.

	logrus.Infof("Aurae runtime daemon. Version: %s", aurae.Version)
	logrus.Infof("Aurae Socket [%s]", d.socket)

	// Step 2. Establish runtime safety

	quitCh := posix.SignalHandler()
	go func() {
		d.runtime = <-quitCh
	}()

	// Step 3. Establish gRPC server.

	var err error
	conn, err := net.Listen("unix", d.socket)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			logrus.Warningf("Attempting to clean socket from failure!")
			err = os.Remove(d.socket)
			if err != nil {
				return fmt.Errorf("unable to establish socket ownership: %v", err)
			}
			conn, err = net.Listen("unix", d.socket)
			if err != nil {
				return fmt.Errorf("unable to establish socket connection: %v", err)
			}
			logrus.Warningf("Socket acquired after previous abandonment.")
		} else {
			return err
		}
	} else {
		logrus.Infof("Success. Socket acquired.")
	}

	server := grpc.NewServer()
	logrus.Infof("Starting gRPC Server.")

	// Step 4. Register the core database to the initialized server

	coreDB := core.NewPathDatabase()
	// TODO We need modular (but opinionated) store backing
	rpc.RegisterCoreServiceServer(server, coreDB)
	logrus.Infof("Registering Core Database.")

	// Step 5. Begin the empty loop by running a small go routine with an emergency cancel

	serveCancel := make(chan error)
	go func() {
		err = server.Serve(conn)
		if err != nil {
			serveCancel <- err
		}
	}()

	logrus.Infof("Listening.")

	// Step 6. Dispatch events from the filesystem

	// need to configure events directories/paths
	//err := d.Dispatch()
	//if err != nil {
	//	logrus.Errorf("Dispatch failure: %v", err)
	//	logrus.Errorf("Shutting down.")
	//	d.runtime = false
	//}

	// Step 7. Begin the runtime loop.

	for d.runtime {
		select {
		case err := <-serveCancel:
			if err != nil {
				logrus.Errorf("Auarae core serving error: %v", err)
				d.runtime = false // Cancel the runtime during a core serving error
			}
		default:
			// pass
		}
	}

	logrus.Infof("Gracefully stopping.")
	server.GracefulStop()
	logrus.Infof("Gentle cleanup.")
	server.Stop()
	logrus.Infof("Safely exiting.")
	return nil
}
