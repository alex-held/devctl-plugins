// +build debug

package plugins

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
)

var Debug_Output_Buffer *bytes.Buffer

func DefaultLoggerOptions() *hclog.LoggerOptions {
	var out io.Writer = os.Stdout
	if Debug_Output_Buffer != nil {
		out = io.MultiWriter(Debug_Output_Buffer, os.Stderr)
	}
	return &hclog.LoggerOptions{
		Name:   "plugin-loader",
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

func log(name string, fn func() error) error {
	start := time.Now()
	defer fmt.Println(name, time.Now().Sub(start))
	return fn()
}
