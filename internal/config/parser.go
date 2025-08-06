package config

import (
	"spooky/internal/config/loading"
	"spooky/internal/config/types"
)

// LoadActionsConfig loads actions configuration from a project
func LoadActionsConfig(projectPath string) (*types.ActionsConfig, error) {
	return loading.LoadActionsConfig(projectPath)
}

// ParseInventoryConfig loads inventory configuration from a file
func ParseInventoryConfig(inventoryFile string) (*types.InventoryConfig, error) {
	return loading.ParseInventoryConfig(inventoryFile)
}

// ParseMachinesInventory loads machines inventory from a file
func ParseMachinesInventory(machinesFile string) (*types.InventoryConfig, error) {
	return loading.ParseMachinesInventory(machinesFile)
}
