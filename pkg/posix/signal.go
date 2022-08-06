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

package posix

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func SignalHandler() chan bool {
	logrus.Debugf("Initializing Signal Handler.")
	sigCh := make(chan os.Signal, 2)
	quitCh := make(chan bool)

	// Register signals for the signal handler
	// os.Interrupt is ^C
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Interrupt)
	go func() {
		sig := <-sigCh
		logrus.Warningf("Shutting down...")
		switch sig {
		case syscall.SIGHUP:
			logrus.Warningf("SIGHUP...")
			quitCh <- false
		case syscall.SIGINT:
			logrus.Warningf("SIGINT...")
			quitCh <- false
		case syscall.SIGTERM:
			logrus.Warningf("SIGTERM...")
			quitCh <- false
		case syscall.SIGKILL:
			logrus.Warningf("SIGKILL...")
			quitCh <- false
		case syscall.SIGQUIT:
			logrus.Warningf("SIGQUIT...")
			quitCh <- false
		default:
			logrus.Warningf("Signal Caught...")
			quitCh <- false
		}
	}()

	return quitCh
}
