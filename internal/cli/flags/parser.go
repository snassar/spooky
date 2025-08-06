package flags

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Parser implements FlagsParser interface
type Parser struct {
	flags map[string]interface{}
}

// NewParser creates a new flags parser
func NewParser() *Parser {
	return &Parser{
		flags: make(map[string]interface{}),
	}
}

// ParseFlags parses flags from a cobra command
func (p *Parser) ParseFlags(cmd interface{}) error {
	if cmd == nil {
		return fmt.Errorf("command cannot be nil")
	}

	// Type assert to cobra.Command
	cobraCmd, ok := cmd.(*cobra.Command)
	if !ok {
		return fmt.Errorf("expected *cobra.Command, got %T", cmd)
	}

	// Parse flags from the command using a simpler approach
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed {
			value := p.getFlagValue(cobraCmd, flag.Name)
			p.flags[flag.Name] = value
		}
	})

	// Validate flags
	if err := p.ValidateFlags(p.flags); err != nil {
		return fmt.Errorf("flag validation failed: %w", err)
	}

	return nil
}

// ValidateFlags validates flags
func (p *Parser) ValidateFlags(flags map[string]interface{}) error {
	if flags == nil {
		return nil
	}

	for name, value := range flags {
		if name == "" {
			return fmt.Errorf("flag name cannot be empty")
		}

		if value == nil {
			return fmt.Errorf("flag value cannot be nil for flag: %s", name)
		}

		// Check if value is of a supported type
		if !p.isSupportedType(value) {
			return fmt.Errorf("unsupported flag type for %s: %T", name, value)
		}
	}

	return nil
}

// GetFlagValue gets a flag value
func (p *Parser) GetFlagValue(name string) interface{} {
	if value, exists := p.flags[name]; exists {
		return value
	}
	return nil
}

// GetStringFlag gets a string flag value
func (p *Parser) GetStringFlag(name string) (string, error) {
	value := p.GetFlagValue(name)
	if value == nil {
		return "", fmt.Errorf("flag %s not found", name)
	}

	if str, ok := value.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("flag %s is not a string", name)
}

// GetIntFlag gets an int flag value
func (p *Parser) GetIntFlag(name string) (int, error) {
	value := p.GetFlagValue(name)
	if value == nil {
		return 0, fmt.Errorf("flag %s not found", name)
	}

	if i, ok := value.(int); ok {
		return i, nil
	}

	return 0, fmt.Errorf("flag %s is not an int", name)
}

// GetBoolFlag gets a bool flag value
func (p *Parser) GetBoolFlag(name string) (bool, error) {
	value := p.GetFlagValue(name)
	if value == nil {
		return false, fmt.Errorf("flag %s not found", name)
	}

	if b, ok := value.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("flag %s is not a bool", name)
}

// GetStringSliceFlag gets a string slice flag value
func (p *Parser) GetStringSliceFlag(name string) ([]string, error) {
	value := p.GetFlagValue(name)
	if value == nil {
		return nil, fmt.Errorf("flag %s not found", name)
	}

	if slice, ok := value.([]string); ok {
		return slice, nil
	}

	return nil, fmt.Errorf("flag %s is not a string slice", name)
}

// Helper method to get flag value from cobra command
func (p *Parser) getFlagValue(cmd *cobra.Command, name string) interface{} {
	flag := cmd.Flags().Lookup(name)
	if flag == nil {
		return nil
	}

	switch flag.Value.Type() {
	case "string":
		return flag.Value.String()
	case "int":
		if i, err := strconv.Atoi(flag.Value.String()); err == nil {
			return i
		}
	case "bool":
		if b, err := strconv.ParseBool(flag.Value.String()); err == nil {
			return b
		}
	case "stringSlice":
		return strings.Split(flag.Value.String(), ",")
	}

	return flag.Value.String()
}

// Helper method to check if a value is of a supported type
func (p *Parser) isSupportedType(value interface{}) bool {
	switch value.(type) {
	case string, int, int64, float64, bool, []string:
		return true
	default:
		return false
	}
}
