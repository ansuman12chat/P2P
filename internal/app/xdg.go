package app

import (
	stdxdg "github.com/adrg/xdg"
)

// The specification defines a set of standard paths for storing application files,
// including data and configuration files. For portability and flexibility reasons,
// applications should use the XDG defined locations instead of hardcoding paths.
type Xdger interface {
	ConfigFile(relPath string) (string, error)
}

type Xdg struct{}

func (a Xdg) ConfigFile(relPath string) (string, error) {
	return stdxdg.ConfigFile(relPath)
}
