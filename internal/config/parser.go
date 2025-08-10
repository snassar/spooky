package config

import (
	spookyconfigloading "spooky/internal/config/loading"
	spookyconfigtypes "spooky/internal/types/config"
)

// LoadActionsConfig loads actions configuration from a project
func LoadActionsConfig(projectPath string) (*spookyconfigtypes.ActionsConfig, error) {
	return spookyconfigloading.LoadActionsConfig(projectPath)
}

// ParseInventoryConfig loads inventory configuration from a file
func ParseInventoryConfig(inventoryFile string) (*spookyconfigtypes.InventoryConfig, error) {
	return spookyconfigloading.ParseInventoryConfig(inventoryFile)
}

// ParseMachinesInventory loads machines inventory from a file
func ParseMachinesInventory(machinesFile string) (*spookyconfigtypes.InventoryConfig, error) {
	return spookyconfigloading.ParseMachinesInventory(machinesFile)
}
