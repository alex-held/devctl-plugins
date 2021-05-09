// +build debug

package devctlog

import (
	"bytes"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

var Debug_Output_Buffer *bytes.Buffer

func DefaultLoggerOptions(name string) *hclog.LoggerOptions {
	var out io.Writer = os.Stdout
	if Debug_Output_Buffer != nil {
		out = io.MultiWriter(Debug_Output_Buffer, os.Stderr)
	}
	return &hclog.LoggerOptions{
		Name:   name,
		Level:  hclog.Trace,
		Output: out,
		// JSONFormat:        false,
		IncludeLocation: true,
		//	TimeFormat:        "",
		DisableTime: true,
		Color:       hclog.ColorOff,
		//	Exclude:           nil,
		// IndependentLevels: false,
	}
}
