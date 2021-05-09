package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/envy"
)

var ErrPluginMissing = fmt.Errorf("plugin is missing")

// LookPath for plugin.
func LookPath(s string) (string, error) {
	if _, err := os.Stat(s); err == nil {
		return s, nil
	}

	if lp, err := exec.LookPath(s); err == nil {
		return lp, err
	}

	var bin string
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var looks []string
	if from, err := envy.MustGet("DEVCTL_PLUGIN_ROOT"); err == nil {
		looks = append(looks, from)
	} else {
		looks = []string{filepath.Join(pwd, "plugins"), filepath.Join(envy.GoPath(), "bin"), envy.Get("PATH", "")}
	}

	for _, p := range looks {
		lp := filepath.Join(p, s)
		if lp, err = filepath.EvalSymlinks(lp); err == nil {
			bin = lp
			break
		}
	}

	if len(bin) == 0 {
		return "", ErrPluginMissing
	}
	return bin, nil
}
