package interfaces

// MachinesIntegration defines the interface for machines system integration
type MachinesIntegration interface {
	// LoadMachines loads machines from the project
	LoadMachines(projectPath string) (*MachinesContext, error)

	// ValidateMachines validates machines data
	ValidateMachines(machines *MachinesContext) error

	// ConnectToMachine connects to a specific machine
	ConnectToMachine(machine string, context *ConnectionContext) error

	// PingMachine pings a machine to check connectivity
	PingMachine(machine string) error

	// GetMachine gets a specific machine by name
	GetMachine(name string, context *MachinesContext) (*Machine, error)

	// ListMachines lists all available machines
	ListMachines(context *MachinesContext) ([]*Machine, error)

	// AddMachine adds a new machine to the project
	AddMachine(name string, machine *Machine, context *MachinesContext) error

	// RemoveMachine removes a machine from the project
	RemoveMachine(name string, context *MachinesContext) error
}

// MachinesContext provides machines data for integrations
type MachinesContext struct {
	BaseContext
	Machines map[string]*Machine
}

// Machine represents a machine in the system
type Machine struct {
	Name     string
	Host     string
	Port     int
	Username string
	Tags     []string
	Metadata map[string]interface{}
}

// ConnectionContext provides context for machine connections
type ConnectionContext struct {
	BaseContext
	Machine *Machine
	Timeout int
	Retries int
}
