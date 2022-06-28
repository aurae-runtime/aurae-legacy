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
	"github.com/kris-nova/aurae/pkg/core/local"
	"github.com/kris-nova/aurae/pkg/posix"
	"github.com/kris-nova/aurae/rpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
)

const (
	DefaultSocketLocationLinux     string = "/run/aurae.sock"
	DefaultLocalStateLocationLinux string = "/var/aurae"
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
	socket     string
	localStore string
}

func New(socket, localStore string) *Daemon {
	return &Daemon{
		runtime:    true,
		socket:     socket,
		localStore: localStore,
	}
}

func (d *Daemon) Run() error {

	// Step 1. Establish context in the logs.

	logrus.Infof("Aurae runtime daemon. Version: %s", aurae.Version)
	logrus.Infof("Aurae Socket [%s]", d.socket)
	logrus.Infof("Aurae Local  [%s]", d.localStore)

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

	// Step 4. Setup the local persistent state.
	localStateStore := local.NewState(d.localStore)
	coreSvc := core.NewService(localStateStore)

	// Default to getFromMemory=false we can change this later
	coreSvc.SetGetFromMemory(false)

	// Step 5. Register the core database to the initialized server
	rpc.RegisterCoreServiceServer(server, coreSvc)
	logrus.Infof("Registering Core Database.")

	// Step 6. Begin the empty loop by running a small go routine with an emergency cancel
	serveCancel := make(chan error)
	go func() {
		err = server.Serve(conn)
		if err != nil {
			serveCancel <- err
		}
	}()

	logrus.Infof("Listening.")

	// Step 7. Dispatch events from the filesystem

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
