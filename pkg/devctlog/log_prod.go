package plugins

import (
	"github.com/hashicorp/go-hclog"
)

func New(name string) hclog.Logger {
	opt := DefaultLoggerOptions()
	opt.Name = name
	return hclog.New(opt)
}

func log(_ string, fn func() error) error {
	return fn()
}
