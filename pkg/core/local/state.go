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

package local

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultStateDirPermission  os.FileMode = 0755
	DefaultStateFilePermission os.FileMode = 0644
)

type State struct {

	// base is the base directory to persist data locally.
	//
	// base will be filepath.Join()'ed to whatever key is passed.
	base string
}

func NewState(base string) *State {
	return &State{
		base: base,
	}
}

func (s *State) Get(key string) string {
	key = filepath.Join(s.base, key)
	data, err := os.ReadFile(key)
	if err != nil {
		logrus.Errorf("Get %s: Failure:  %v", key, err)
	} else {
		logrus.Debugf("Get %s: Success", key)
	}
	return string(data)
}

func (s *State) Set(key, value string) {
	if strings.HasSuffix(key, "/") {
		key = filepath.Join(s.base, key)
		key = fmt.Sprintf("%s/", key)
	} else {
		key = filepath.Join(s.base, key)
	}
	err := os.MkdirAll(filepath.Dir(key), DefaultStateDirPermission)
	if err != nil {
		logrus.Errorf("Set %s: Failure:  %v", key, err)
		return
	}
	err = os.WriteFile(key, []byte(value), DefaultStateFilePermission)
	if err != nil {
		logrus.Errorf("Set %s: Failure:  %v", key, err)
	} else {
		logrus.Debugf("Set %s: Success", key)
	}
}

func (s *State) List(key string) map[string]bool {
	key = filepath.Join(s.base, key)
	dirents, err := os.ReadDir(key)
	if err != nil {
		logrus.Errorf("List %s: Failure:  %v", key, err)
	} else {
		logrus.Debugf("List %s: Success", key)
	}
	ret := make(map[string]bool)
	for _, dirent := range dirents {
		ret[dirent.Name()] = !dirent.IsDir()
	}
	return ret
}

func (s *State) Remove(key string) {
	key = filepath.Join(s.base, key)
	err := os.RemoveAll(key)
	if err != nil {
		logrus.Errorf("Remove %s: Failure:  %v", key, err)
	} else {
		logrus.Debugf("Remove %s: Success", key)
	}
}
