package config

const (
	// DefaultSSHPort is the default SSH port if not specified
	DefaultSSHPort = 22

	// DefaultTimeout is the default timeout for SSH connections in seconds
	DefaultTimeout = 30

	// DefaultPasswordLength is the default length for generated passwords
	DefaultPasswordLength = 25

	// MaxKeyDirectories is the maximum number of key directories per day
	MaxKeyDirectories = 1000
)

// SetMachineDefaults applies default values to a machine
func SetMachineDefaults(machine *Machine) {
	if machine == nil {
		return
	}
	if machine.Port == 0 {
		machine.Port = DefaultSSHPort
	}
}

// SetActionDefaults applies default values to an action
func SetActionDefaults(action *Action) {
	if action == nil {
		return
	}
	if action.Timeout == 0 {
		action.Timeout = DefaultTimeout
	}
}

// SetInventoryDefaults applies default values to an inventory configuration
func SetInventoryDefaults(inventory *InventoryConfig) {
	if inventory == nil {
		return
	}
	for i := range inventory.Machines {
		SetMachineDefaults(&inventory.Machines[i])
	}
}

// SetActionsDefaults applies default values to an actions configuration
func SetActionsDefaults(actions *ActionsConfig) {
	if actions == nil {
		return
	}
	for i := range actions.Actions {
		SetActionDefaults(&actions.Actions[i])
	}
}
