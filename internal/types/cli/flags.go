package cli

// Flag represents a CLI flag
type Flag struct {
	Name       string      `hcl:"name"`
	Short      string      `hcl:"short,optional"`
	Usage      string      `hcl:"usage"`
	Default    interface{} `hcl:"default,optional"`
	Required   bool        `hcl:"required,optional"`
	Persistent bool        `hcl:"persistent,optional"`
	Hidden     bool        `hcl:"hidden,optional"`
	Deprecated bool        `hcl:"deprecated,optional"`
	Shorthand  string      `hcl:"shorthand,optional"`
	ValueType  string      `hcl:"value_type,optional"`
}

// FlagSet represents a set of flags
type FlagSet struct {
	Global  []*Flag            `hcl:"global,optional"`
	Command map[string][]*Flag `hcl:"command,optional"`
}

// FlagValue represents a flag value
type FlagValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}
