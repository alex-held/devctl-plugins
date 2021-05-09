// +build !debug

package plugins

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

func DefaultLoggerOptions() *hclog.LoggerOptions {
	return &hclog.LoggerOptions{
		Name:   "plugin-loader",
		Level:  hclog.Info,
		Output: os.Stdout,
		// JSONFormat:        false,
		IncludeLocation: true,
		//	TimeFormat:        "",
		DisableTime: true,
		Color:       hclog.AutoColor,
		//	Exclude:           nil,
		// IndependentLevels: false,
	}
}

// Logger the
var Logger = hclog.New(DefaultLoggerOptions())

func New(name string) hclog.Logger {
	opt := DefaultLoggerOptions()
	opt.Name = name
	return hclog.New(opt)
}

func log(_ string, fn func() error) error {
	return fn()
}
