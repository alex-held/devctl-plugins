package deps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"

	"github.com/alex-held/devctl-plugins/pkg/devctlog"
	"github.com/alex-held/devctl-plugins/pkg/devctlpath"
)

// List all of the plugins the application depends on. Will return ErrMissingConfig
// if the app is not using $DEVCTL_ROOT/plugins.toml to manage their plugins.
// Use plugdeps#On(app) to test if plugdeps are being used.
func List(logger devctlog.Logger, pather devctlpath.Pather) (*Plugins, error) {
	plugs := New()
	lp, err := listLocal(logger, pather)
	if err != nil {
		return plugs, err
	}
	plugs.Add(lp.List()...)

	p := pather.PluginConfigFilePath()
	tf, err := os.Open(p)
	if err != nil {
		return plugs, err
	}

	if err := plugs.Decode(tf); err != nil {
		logger.Error("decoding of $DEVCTL_ROOT/plugins.toml failed", "err", err, "file", *tf, "plugs", plugs)
		return plugs, err
	}

	local, err := listLocal(logger, pather)
	if err != nil {
		logger.Error("discovering listLocal failed", "err", err, "plugs", plugs)
		return plugs, err
	}
	plugs.Add(local.List()...)

	return plugs, nil
}

func listLocal(logger devctlog.Logger, pather devctlpath.Pather) (*Plugins, error) {
	plugs := New()
	origLoggerName := logger.Name()

	proot := pather.Plugins()
	if _, err := os.Stat(proot); err != nil {
		return plugs, nil
	}

	dirwalkLogger := logger.ResetNamed(fmt.Sprintf("%s | godirwalk\t", origLoggerName))
	err := godirwalk.Walk(proot, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback: func(path string, info *godirwalk.Dirent) error {
			dirwalkLogger.Debug(fmt.Sprintf("walking in '%s'", path), "skipping", info.IsDir(), "fileInfo", *info, "fi_name", info.Name(), "plugin_root", proot)

			if info.IsDir() {
				return nil
			}

			base := filepath.Base(path)
			if strings.HasPrefix(base, "devctl-") {
				local := "." + strings.TrimPrefix(path, pather.ConfigRoot())
				plugs.Add(Plugin{
					Binary: base,
					Local:  local,
				})
			}
			return nil
		},
	})
	if err != nil {
		dirwalkLogger.Error("resolving plugins via dir-walking failed", "err", err)
		return plugs, err
	}

	dirwalkLogger.ResetNamed(origLoggerName)
	return plugs, nil
}
