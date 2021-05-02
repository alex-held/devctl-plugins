package sdk

import (
	"github.com/gobuffalo/plugins/plugcmd"
)

type Installer interface {
	Install(version string) error
	plugcmd.Namer
}

type Downloader interface {
	Install(version string) error
	plugcmd.Namer
}

type Lister interface {
	ListInstalled(version string) ([]string, error)
	plugcmd.Namer
}

type Linker interface {
	Link(version string) error
	plugcmd.Namer
}

type User interface {
	Use(version string) error
	plugcmd.Namer
}
