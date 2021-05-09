package devctlog

import (
	"github.com/hashicorp/go-hclog"
)

type Logger = hclog.Logger

func New(name string) Logger {
	opts := DefaultLoggerOptions(name)
	return hclog.New(opts)
}
