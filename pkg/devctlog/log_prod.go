// +build !debug

package devctlog

import (
	"os"

	"github.com/hashicorp/go-hclog"
)

func DefaultLoggerOptions(name string) *hclog.LoggerOptions {
	return &hclog.LoggerOptions{
		Name:   name,
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
