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
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// Dispatch is the primary event broker. From here we will be able to dispatch
// messages in the system to configured response mechanisms.
func (d *Daemon) Dispatch() error {
	logrus.Infof("Initalizing auraeFS watcher")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("unable to start dispatch watcher: %v", err)
	}
	defer watcher.Close()
	go func() {
		for d.runtime {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {

					// TODO We need a dynamic dispatch mechanism. Ideally something simple.

					logrus.Infof("Modified file: %s", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Warnf("Event broker error: %v", err)
			}
		}
	}()

	// Register the aurae filesystem
	var path string

	// /aurae
	// TODO Make this a modular system we can use dependency injection
	path = "/aurae"
	logrus.Infof("Watching [%s] for events...", path)
	err = watcher.Add(path)
	if err != nil {
		return fmt.Errorf("unable to watch: %s: %v", path, err)
	}

	return nil
}
