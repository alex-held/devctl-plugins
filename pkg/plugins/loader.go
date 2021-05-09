package plugins

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/events"
	"github.com/hashicorp/go-hclog"
	"github.com/markbates/oncer"
	"github.com/markbates/safe"

	"github.com/alex-held/devctl-plugins/pkg/devctlog"
)

// LoadPlugins will add listeners for any plugins that support "events".
func LoadPlugins() error {
	opts := devctlog.DefaultLoggerOptions("events.LoadPlugins")
	opts.Level = hclog.Trace
	opts.TimeFormat = "[15:04:05.00]    "
	opts.DisableTime = false
	opts.JSONFormat = true
	opts.Color = hclog.AutoColor
	opts.IncludeLocation = true
	logger := hclog.New(opts)
	var err error

	oncer.Do("events.LoadPlugins", func() {
		// don't send plugins events during testing
		if envy.Get("GO_ENV", "development") == "test" {
			return
		}

		plugs, err := Available()
		if err != nil {
			err = err
			return
		}

		for _, cmds := range plugs {
			for _, c := range cmds {
				if c.DevctlCommand != "events" {
					continue
				}

				err := func(c Command) error {
					return safe.RunE(func() error {
						n := fmt.Sprintf("[PLUGIN] %s %s", c.Binary, c.Name)
						fn := func(e events.Event) {
							eventBytes, err := json.Marshal(e)
							if err != nil {
								logger.Error("error trying to marshal event", "event", e, "err", err)
								return
							}

							jsonEvent := string(eventBytes)
							logger.Trace("serialized event to json", "event", e, "json", jsonEvent)
							cmd := exec.Command(c.Binary, c.UseCommand, jsonEvent)
							errLogWriter := logger.Named(c.Binary + c.UseCommand + "STD::ERR").StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
							outLogWriter := logger.Named(c.Binary + c.UseCommand + "STD::OUT").StandardWriter(&hclog.StandardLoggerOptions{InferLevels: true})
							cmd.Stderr = io.MultiWriter(os.Stderr, errLogWriter)
							cmd.Stdout = io.MultiWriter(os.Stdout, outLogWriter)
							cmd.Stdin = os.Stdin

							logger.Debug("starting to execute command", "command", *cmd, "HumanReadableCommand", cmd.String())

							if err := cmd.Run(); err != nil {
								logger.Error("error trying to send event", "command", *cmd, "args", strings.Join(cmd.Args, " "), "err", err)
							}
						}
						_, err := events.NamedListen(n, events.Filter(c.ListenFor, fn))
						if err != nil {
							logger.Error("error while listening...", "err", err)
							return err
						}
						return nil
					})
				}(c)

				if err != nil {
					err = err
					return
				}
			}
		}
	})
	return err
}
