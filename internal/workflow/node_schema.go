package workflow

type (
	// NodeSchema represents a structure containing an ID, package details, and optional configuration.
	NodeSchema struct {
		ID       string      `json:"id" yaml:"id"`
		Function string      `json:"function" yaml:"function"`
		Config   *NodeConfig `json:"config,omitempty" yaml:"config,omitempty"`
	}
	// NodeConfig represents the configuration schema for a node. TODO
	NodeConfig struct{}
)
