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

package daemon

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae"
	"github.com/kris-nova/aurae/pkg/core"
	"github.com/kris-nova/aurae/pkg/core/local"
	p2pgrpc "github.com/kris-nova/aurae/pkg/grpc"
	"github.com/kris-nova/aurae/pkg/name"
	"github.com/kris-nova/aurae/pkg/peer"
	"github.com/kris-nova/aurae/pkg/posix"
	"github.com/kris-nova/aurae/pkg/proxy"
	"github.com/kris-nova/aurae/pkg/runtime"
	"github.com/kris-nova/aurae/pkg/schedule"
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

	// Self is the root of the peer to peer digraph.
	Self *peer.Peer

	// Config options
	runtime    bool
	socket     string
	localStore string
	//keypath    string
}

func New(socket, localStore string) *Daemon {
	return &Daemon{
		runtime:    true,
		socket:     socket,
		localStore: localStore,
		//keypath:    keypath,
	}
}

func (d *Daemon) Run() error {

	ctx := context.Background()

	// Establish context in the logs.
	logrus.Infof("Aurae daemon daemon. Version: %s", aurae.Version)
	logrus.Infof("Aurae Socket [%s]", d.socket)
	logrus.Infof("Aurae Local  [%s]", d.localStore)

	// Establish daemon safety
	quitCh := posix.SignalHandler()
	go func() {
		d.runtime = <-quitCh
	}()

	// Establish gRPC server.
	var err error
	socketConn, err := net.Listen("unix", d.socket)
	if err != nil {
		if strings.Contains(err.Error(), "bind: address already in use") {
			logrus.Warningf("Attempting to clean socket from failure!")
			err = os.Remove(d.socket)
			if err != nil {
				return fmt.Errorf("unable to establish socket ownership: %v", err)
			}
			socketConn, err = net.Listen("unix", d.socket)
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
	defer socketConn.Close()

	// Local gRPC server
	logrus.Infof("Starting [SOCKET] gRPC Server.")
	server := grpc.NewServer()

	// Step 4. Setup the local persistent state.
	localStateStore := local.NewState(d.localStore)
	coreSvc := core.NewService(localStateStore)

	// Default to getFromMemory=false we can change this later
	coreSvc.SetGetFromMemory(false)

	//
	rpc.RegisterCoreServer(server, coreSvc)
	rpc.RegisterProxyServer(server, proxy.NewService())
	rpc.RegisterRuntimeServer(server, runtime.NewService())
	rpc.RegisterScheduleServer(server, schedule.NewService())

	logrus.Infof("Registering Core Services.")

	// Step 6. Begin the empty loop by running a small go routine with an emergency cancel
	serveCancel := make(chan error)
	go func() {
		err = server.Serve(socketConn)
		if err != nil {
			serveCancel <- err
		}
	}()

	// Step 7. Peer host and initialize peer to peer network.
	//instanceKey, err := crypto.KeyFromPath(d.keypath)
	//if err != nil {
	//	return fmt.Errorf("invalid private key: %s: %v", d.keypath, err)
	//}
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("unable to calculate hostname: %v", err)
	}
	self := peer.NewPeer(name.New(hostname))
	err = self.Establish(context.Background(), 1)
	if err != nil {
		return fmt.Errorf("unable to join auraespace peer network: %v", err)
	}
	go self.HandshakeServe()
	d.Self = self
	logrus.Infof("Starting Auare handshake protocol on peer network")

	peerConn := p2pgrpc.NewGRPCProtocol(ctx, self.Host())
	if err != nil {
		return fmt.Errorf("unable to create peer grpc: %v", err)
	}
	server = peerConn.GetGRPCServer()
	rpc.RegisterCoreServer(server, coreSvc)
	rpc.RegisterProxyServer(server, proxy.NewService())
	rpc.RegisterRuntimeServer(server, runtime.NewService())
	rpc.RegisterScheduleServer(server, schedule.NewService())
	logrus.Infof("Starting Auare grpc protocol on peer network")

	//logrus.Infof("Setting peer grpc: %v", server.GetServiceInfo())

	// Step 7. Dispatch events from the filesystem

	// need to configure events directories/paths
	//err := d.Dispatch()
	//if err != nil {
	//	logrus.Errorf("Dispatch failure: %v", err)
	//	logrus.Errorf("Shutting down.")
	//	d.daemon = false
	//}

	// Step 7. Begin the daemon loop.

	for d.runtime {
		select {
		case err := <-serveCancel:
			if err != nil {
				logrus.Errorf("Auarae core serving error: %v", err)
				d.runtime = false // Cancel the daemon during a core serving error
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
