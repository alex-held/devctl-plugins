package sdk

import (
	"context"

	"github.com/gobuffalo/plugins/plugcmd"
)

type Installer interface {
	plugcmd.Namer
	Install(ctx context.Context, version string) error
}

type Downloader interface {
	plugcmd.Namer
	Download(ctx context.Context, version string) error
}

type Lister interface {
	plugcmd.Namer
	ListInstalled(ctx context.Context) ([]string, error)
}

type Linker interface {
	plugcmd.Namer
	Link(ctx context.Context, version string) error
}

type User interface {
	plugcmd.Namer
	Use(ctx context.Context, version string) error
}
