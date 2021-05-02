package plugins

import (
	"github.com/gobuffalo/plugins"
	"github.com/gobuffalo/plugins/plugcmd"

	"github.com/alex-held/devctl-plugins/sdk"
)

type SDKPlugin interface {
	plugcmd.Commander

	plugins.Needer
	plugcmd.SubCommander

	sdk.Downloader
	sdk.Lister
	sdk.Linker
	sdk.Installer
	sdk.User
}
