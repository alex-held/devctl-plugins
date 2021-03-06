package plugins

// Commands is a slice of Command.
type Commands []Command

// Command that the plugin supplies.
type Command struct {
	// Name "foo"
	Name string `json:"name"`

	// UseCommand "bar"
	UseCommand string `json:"use_command"`

	// DevctlCommand "generate"
	DevctlCommand string `json:"devctl_command"`

	// Description "generates a foo"
	Description string   `json:"description,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	Binary      string   `json:"-"`
	Flags       []string `json:"flags,omitempty"`

	// Filters events to listen to ("" or "*") is all events
	ListenFor string `json:"listen_for,omitempty"`
}
