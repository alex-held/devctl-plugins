package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/alex-held/devctl-plugins/pkg/devctlog"
	"github.com/alex-held/devctl-plugins/pkg/sysutils"
)

var _logger devctlog.Logger

func log(name string, fn func() error) error {
	start := time.Now()
	defer _logger.Debug(name, time.Now().Sub(start))
	return fn()
}

// Decorate setup cmds Commands for plugins.
func Decorate(c Command) *cobra.Command {
	var flags []string
	if len(c.Flags) > 0 {
		for _, f := range c.Flags {
			flags = append(flags, f)
		}
	}

	cmd := &cobra.Command{
		Use:     c.Name,
		Short:   fmt.Sprintf("[PLUGIN] %s", c.Description),
		Aliases: c.Aliases,
		RunE: func(runCmd *cobra.Command, args []string) error {
			_logger.Debug("executing`devctl-plugin/pkg/plugins/Plugin` decorated as a `cmds.Command`",
				"cmd", *runCmd, "args", args)

			plugCmd := c.Name
			if c.UseCommand != "" {
				plugCmd = c.UseCommand
			}

			ax := []string{plugCmd}
			if plugCmd == "-" {
				ax = []string{}
			}

			ax = append(ax, args...)
			ax = append(ax, flags...)

			bin, err := LookPath(c.Binary)
			_logger.Debug("LookPath finished..", "looked-up-binary", c.Binary, "bin", bin, "err", err)
			if err != nil {
				_logger.Error("LookPath failed to look up binary", "looked-up-binary", c.Binary, "bin", bin, "err", err)
				return err
			}

			ex := exec.Command(bin, ax...)
			if !sysutils.GetDefaultRuntimeInfo().IsWindows() {
				ex.Env = append(os.Environ(), "DEVCTL_PLUGIN=1")
			}

			ex.Stdin = os.Stdin
			ex.Stdout = os.Stdout
			ex.Stderr = os.Stderr

			return log(c.Binary, ex.Run)
		},
	}
	cmd.DisableFlagParsing = true
	return cmd
}
