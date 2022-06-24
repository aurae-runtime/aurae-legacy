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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/shape"
)

var _ shape.AuraeFilesystem = &LocalFilesystem{}

// LocalFilesystem is a simple local state store for aurae
type LocalFilesystem struct {
	root  string
	model *shape.Model
}

// NewFromPath assumes that an shape is in place, otherwise will error.
//
// New from path will validate the path does not exist, and create an empty
// root for shape.
//
// We call this function out so it may be overloaded later if we decide to
// implement a virtual filesystem.
func NewFromPath(path string) (*LocalFilesystem, error) {
	path = common.Expand(path)
	if !strings.Contains(path, "aurae") {
		path = filepath.Join(path, "aurae")
	}
	_, err := os.Stat(path)
	if err == nil {
		return nil, fmt.Errorf("%s exists", path)
	}
	// The directory does not exist, initialize a new root for the new shape
	err = os.MkdirAll(path, shape.DefaultMode)
	if err != nil {
		return nil, fmt.Errorf("unable to create new shape root: %s: %v", path, err)
	}
	return &LocalFilesystem{
		root:  path,
		model: shape.NewModel(path),
	}, nil
}

// LoadFromPath assumes that an shape is in place, otherwise will error
func LoadFromPath(path string) (*LocalFilesystem, error) {
	path = common.Expand(path)
	if !strings.Contains(path, "aurae") {
		path = filepath.Join(path, "aurae")
	}
	s, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load shape from path %s: %v", path, err)
	}
	if !s.IsDir() {
		return nil, fmt.Errorf("unable to load shape from path %s: not directory", path)
	}
	return &LocalFilesystem{
		root:  path,
		model: shape.NewModel(path),
	}, nil
}

// Load assumes an shape has already been initialized, and will look in common locations
func Load() (*LocalFilesystem, error) {
	pathsToTry := []string{
		"/aurae",
		"~aurae",
		"/usr/local/aurae",
		"/var/run/aurae",
	}
	var lfs *LocalFilesystem
	var err error
	for _, pathToTry := range pathsToTry {
		lfs, err = LoadFromPath(pathToTry)
		if err != nil {
			continue
		}
		break
	}
	if lfs == nil {
		return nil, fmt.Errorf("unable to find shape from common paths. tried paths: %s", strings.Join(pathsToTry, ","))
	}
	return lfs, nil
}

func (a *LocalFilesystem) Destroy() error {

	// Here we safeguard against a classic rm -rf /
	if a.root == "/" {
		return fmt.Errorf("unable to destroy / for aurfs")
	}

	return os.RemoveAll(a.root)
}

func (a *LocalFilesystem) Initialize(cfg *shape.Config) error {

	// Directories first
	for dir, mode := range a.model.Directories {
		err := os.MkdirAll(filepath.Join(a.root, dir), mode)
		if err != nil {
			return err
		}
	}

	// Files second
	for file, mode := range a.model.Files {
		err := os.WriteFile(filepath.Join(a.root, file), shape.Empty, mode)
		if err != nil {
			return err
		}
	}

	// Now process the config
	err := a.Set(shape.FileEtcNameNode, cfg.NodeName)
	if err != nil {
		return err
	}
	err = a.Set(shape.FileEtcNameDomain, cfg.DomainName)
	if err != nil {
		return err
	}
	return nil
}

func (a *LocalFilesystem) MountPoint() string {
	return a.root
}

func (a *LocalFilesystem) Socket() string {
	return filepath.Join(a.root, shape.FileauraeSock)
}

func (a *LocalFilesystem) Set(key, value string) error {
	err := os.MkdirAll(path.Dir(filepath.Join(a.root, key)), shape.DefaultMode)
	if err != nil {
		return fmt.Errorf("set() error: %v", err)
	}
	return os.WriteFile(filepath.Join(a.root, key), []byte(value), shape.DefaultMode)
}

func (a *LocalFilesystem) Get(key string) string {
	data, err := os.ReadFile(filepath.Join(a.root, key))
	if err != nil {
		return ""
	}
	return string(data)
}
