package plugins

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"

	sdk2 "github.com/alex-held/devctl-plugins/pkg/plugins/sdk"
)

type SDKPlugin interface {
	plugcmd.Commander
	Versioner

	plugins.Needer
	plugcmd.SubCommander

	sdk2.Downloader
	sdk2.Lister
	sdk2.Linker
	sdk2.Installer
	sdk2.User
}
