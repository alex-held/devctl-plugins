package deps

import (
	"encoding/json"
	"fmt"

	"github.com/gobuffalo/meta"
)

// ErrMissingConfig is if $DEVCTL_ROOT/plugins.toml file is not found. Use plugdeps#On(app) to test if plugdeps are being used.
var ErrMissingConfig = fmt.Errorf("could not find a plugin configuration at `$DEVCTL_ROOT/plugins.toml`")

// Plugin represents a Go plugin for Buffalo applications.
type Plugin struct {
	Binary   string         `toml:"binary" json:"binary"`
	GoGet    string         `toml:"go_get,omitempty" json:"go_get,omitempty"`
	Local    string         `toml:"local,omitempty" json:"local,omitempty"`
	Commands []Command      `toml:"commands,omitempty" json:"commands,omitempty"`
	Tags     meta.BuildTags `toml:"tags,omitempty" json:"tags,omitempty"`
}

// String implementation of fmt.Stringer.
func (p Plugin) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func (p Plugin) key() string {
	return p.Binary + p.GoGet + p.Local
}
