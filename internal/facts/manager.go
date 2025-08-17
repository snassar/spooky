// Package facts provides fact collection and in-memory management functionality.
package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookyschemas "spooky/internal/schemas"
	spookysecrets "spooky/internal/secrets"
	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesschemas "spooky/internal/types/schemas"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Manager implements FactManager
type Manager struct {
	collector             spookytypes.FactCollector
	schemaDrivenValidator *spookyschemas.SchemaDrivenValidator
	enhancedValidator     *spookyschemas.EnhancedValidator
	logger                spookytypeslogging.Logger
}

// NewManager creates a new fact manager
func NewManager(
	collector spookytypes.FactCollector,
	logger spookytypeslogging.Logger,
) *Manager {
	// Create schema-driven validator for fact configuration validation
	schemaDrivenConfig := &spookyschemas.SchemaDrivenValidationConfig{
		UseEmbeddedSchemas: true,
		StrictValidation:   true,
		AllowUnknownFields: false,
		DetailedErrors:     true,
	}
	schemaDrivenValidator := spookyschemas.NewSchemaDrivenValidator(logger, schemaDrivenConfig)

	// Create enhanced validator for fact data validation
	enhancedConfig := &spookyschemas.ValidationConfig{
		Mode: spookyschemas.ValidationModeStrict,
		ErrorHandling: &spookyschemas.ErrorHandlingConfig{
			StopOnFirstError:   false,
			MaxErrors:          100,
			IncludeWarnings:    true,
			IncludeContext:     true,
			IncludeSuggestions: true,
		},
		Evolution: &spookyschemas.EvolutionConfig{
			EnableTracking:  true,
			AllowDeprecated: true,
			WarnDeprecated:  true,
			AllowBreaking:   false,
		},
	}
	enhancedValidator := spookyschemas.NewEnhancedValidator(enhancedConfig)

	return &Manager{
		collector:             collector,
		schemaDrivenValidator: schemaDrivenValidator,
		enhancedValidator:     enhancedValidator,
		logger:                logger,
	}
}

// CollectFacts collects facts from the given machine
func (m *Manager) CollectFacts(ctx context.Context, machine interface{}) (*spookytypes.FactCollection, error) {
	if machine == nil {
		return nil, fmt.Errorf("machine cannot be nil")
	}

	// Type assert to get machine details
	machineObj, ok := machine.(*spookytypes.Machine)
	if !ok {
		return nil, fmt.Errorf("machine must be of type *spookytypes.Machine")
	}

	m.logger.Info("Collecting facts", map[string]interface{}{
		"machine": machineObj.Hostname,
		"host":    machineObj.Host,
	})

	// Determine collection method based on machine configuration
	if machineObj.Host != "" && machineObj.Host != "localhost" && machineObj.Host != "127.0.0.1" {
		// Remote machine - use SSH-based collection
		return m.collectFactsViaSSH(ctx, machineObj)
	}

	// Local machine - use local collection
	facts, err := m.collector.Collect(ctx, machineObj)
	if err != nil {
		m.logger.Error("Failed to collect facts", err, map[string]interface{}{"machine": machineObj.Hostname})
		return nil, fmt.Errorf("failed to collect facts for %s: %w", machineObj.Hostname, err)
	}

	m.logger.Info("Successfully collected facts", map[string]interface{}{
		"machine":   machineObj.Hostname,
		"collector": m.collector.GetName(),
	})

	return facts, nil
}

// collectFactsViaSSH collects facts from remote machine via SSH
func (m *Manager) collectFactsViaSSH(ctx context.Context, machine *spookytypes.Machine) (*spookytypes.FactCollection, error) {
	m.logger.Info("Collecting facts via SSH", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	// Get SSH-capable collector
	sshCollector, ok := m.collector.(interface {
		CollectViaSSH(context.Context, *spookytypes.Machine) (*spookytypes.FactCollection, error)
	})
	if !ok {
		return nil, fmt.Errorf("collector does not support SSH operations")
	}

	// Collect facts using SSH
	facts, err := sshCollector.CollectViaSSH(ctx, machine)
	if err != nil {
		m.logger.Error("Failed to collect facts via SSH", err, map[string]interface{}{
			"machine": machine.Hostname,
			"host":    machine.Host,
		})
		return nil, fmt.Errorf("failed to collect facts via SSH for %s: %w", machine.Hostname, err)
	}

	m.logger.Info("Successfully collected facts via SSH", map[string]interface{}{
		"machine": machine.Hostname,
		"host":    machine.Host,
	})

	return facts, nil
}

// GetCollector returns the underlying collector for SSH operations
func (m *Manager) GetCollector() spookytypes.FactCollector {
	return m.collector
}

// ValidateFacts validates fact collection using schema validators
func (m *Manager) ValidateFacts(ctx context.Context, facts *spookytypes.FactCollection) (interface{}, error) {
	m.logger.Info("Validating facts", map[string]interface{}{
		"machine_id": facts.MachineID,
	})

	// Get facts schema for enhanced validation
	factsSchema, err := m.getFactsSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to get facts schema: %w", err)
	}

	// Use enhanced validator for comprehensive fact validation
	result, err := m.enhancedValidator.ValidateWithEnhancedFeatures(ctx, factsSchema, facts)
	if err != nil {
		return nil, fmt.Errorf("failed to validate facts with enhanced validator: %w", err)
	}

	// Add additional custom validation for fact-specific rules
	m.addCustomFactValidation(facts, result)

	return result, nil
}

// getFactsSchema gets the facts schema for validation
func (m *Manager) getFactsSchema() (*spookytypesschemas.Schema, error) {
	// Try to get schema from embedded schemas first
	if schema, err := m.schemaDrivenValidator.GetEmbeddedSchema("facts"); err == nil {
		return schema, nil
	}

	// Fallback: create a basic facts schema
	return &spookytypesschemas.Schema{
		Name:        "facts",
		Type:        "hcl",
		Version:     "1.0",
		Description: "Facts collection schema",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Content:     "", // Will be loaded from file if needed
		Metadata:    make(map[string]interface{}),
	}, nil
}

// addCustomFactValidation adds custom validation rules specific to facts
func (m *Manager) addCustomFactValidation(facts *spookytypes.FactCollection, result *spookytypesschemas.ValidationResult) {
	// Validate machine ID
	if facts.MachineID == "" {
		m.addSchemaError(result, "missing_machine_id", "Machine ID is required", "error")
	}

	// Validate fact collection timestamp
	if facts.CollectedAt.IsZero() {
		m.addSchemaError(result, "missing_collection_timestamp", "Collection timestamp is required", "error")
	}

	// Validate facts structure
	if facts.Facts == nil {
		m.addSchemaError(result, "missing_facts_structure", "Facts structure is required", "error")
		return
	}

	// Validate system facts
	if facts.Facts.System == nil {
		m.addSchemaError(result, "missing_system_facts", "System facts are required", "error")
	} else {
		// Validate OS facts
		if facts.Facts.System.OS == nil {
			m.addSchemaError(result, "missing_os_facts", "OS facts are required", "error")
		}

		// Validate hardware facts
		if facts.Facts.System.Hardware == nil {
			m.addSchemaError(result, "missing_hardware_facts", "Hardware facts are required", "error")
		}

		// Validate network facts
		if facts.Facts.System.Network == nil {
			m.addSchemaError(result, "missing_network_facts", "Network facts are required", "error")
		}
	}
}

// addSchemaError adds a schema error to the validation result
func (m *Manager) addSchemaError(result *spookytypesschemas.ValidationResult, code, message, severity string) {
	schemaError := spookytypesschemas.SchemaError{
		Code:     code,
		Message:  message,
		Severity: severity,
	}
	result.Errors = append(result.Errors, schemaError)
	result.Valid = false
}

// ExportFacts exports facts to the given format
func (m *Manager) ExportFacts(ctx context.Context, machineIDs []string, format, outputPath string) error {
	m.logger.Info("Exporting facts", map[string]interface{}{
		"machines": len(machineIDs),
		"format":   format,
		"output":   outputPath,
	})

	// Collect facts for the specified machines
	var allFacts []*spookytypes.FactCollection
	for _, machineID := range machineIDs {
		facts, err := m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
		if err != nil {
			m.logger.Error("Failed to collect facts for machine", err, map[string]interface{}{
				"machine": machineID,
			})
			return fmt.Errorf("failed to collect facts for machine %s: %w", machineID, err)
		}
		allFacts = append(allFacts, facts)
	}

	// Export based on format
	switch format {
	case "json":
		return m.exportToJSON(allFacts, outputPath)
	case "hcl":
		return m.exportToHCL(allFacts, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToJSON exports facts to JSON format following facts-structure.schema.hcl
func (m *Manager) exportToJSON(facts []*spookytypes.FactCollection, outputPath string) error {
	// Export each fact collection individually following the schema structure
	var exportData []map[string]interface{}

	for _, fact := range facts {
		if fact == nil {
			continue
		}

		// Create the facts_structure format
		factStructure := map[string]interface{}{
			"machine_id":   fact.MachineID,
			"collected_at": fact.CollectedAt.Format(time.RFC3339),
			"facts":        fact.Facts,
		}

		exportData = append(exportData, factStructure)
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal facts to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write JSON export file: %w", err)
	}

	m.logger.Info("Successfully exported facts to JSON", map[string]interface{}{"output": outputPath, "machines": len(facts)})
	return nil
}

// exportToHCL exports facts to HCL format following facts-structure.schema.hcl
func (m *Manager) exportToHCL(facts []*spookytypes.FactCollection, outputPath string) error {
	file := hclwrite.NewEmptyFile()
	rootBody := file.Body()

	// Export each fact collection following the facts_structure schema
	for idx := range facts {
		fact := facts[idx]
		if fact == nil {
			continue
		}

		m.addFactStructureToHCL(rootBody, fact)
	}

	// Write to file
	if err := os.WriteFile(outputPath, file.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write HCL export file: %w", err)
	}

	m.logger.Info("Successfully exported facts to HCL", map[string]interface{}{"output": outputPath, "machines": len(facts)})
	return nil
}

// addFactStructureToHCL adds a fact structure block to the HCL body
func (m *Manager) addFactStructureToHCL(rootBody *hclwrite.Body, fact *spookytypes.FactCollection) {
	// Create facts_structure block
	factsStructureBlock := rootBody.AppendNewBlock("facts_structure", nil)
	factsStructureBody := factsStructureBlock.Body()

	// Add machine_id
	factsStructureBody.SetAttributeValue("machine_id", cty.StringVal(fact.MachineID))

	// Add collected_at
	factsStructureBody.SetAttributeValue("collected_at", cty.StringVal(fact.CollectedAt.Format(time.RFC3339)))

	// Add facts object
	if fact.Facts != nil {
		factsBlock := factsStructureBody.AppendNewBlock("facts", nil)
		factsBody := factsBlock.Body()

		// Add system facts if available
		if fact.Facts.System != nil {
			m.addSystemFactsToHCL(factsBody, fact.Facts.System)
		}
	}
}

// addSystemFactsToHCL adds system facts to the HCL body
func (m *Manager) addSystemFactsToHCL(factsBody *hclwrite.Body, system *spookytypesfacts.SystemFacts) {
	systemBlock := factsBody.AppendNewBlock("system", nil)
	systemBody := systemBlock.Body()

	// OS facts
	if system.OS != nil {
		m.addOSFactsToHCL(systemBody, system.OS)
	}

	// Hardware facts
	if system.Hardware != nil {
		m.addHardwareFactsToHCL(systemBody, system.Hardware)
	}

	// Network facts
	if system.Network != nil {
		m.addNetworkFactsToHCL(systemBody, system.Network)
	}
}

// addOSFactsToHCL adds OS facts to the HCL body
func (m *Manager) addOSFactsToHCL(systemBody *hclwrite.Body, os *spookytypesfacts.OSFacts) {
	osBlock := systemBody.AppendNewBlock("os", nil)
	osBody := osBlock.Body()
	osBody.SetAttributeValue("name", cty.StringVal(os.Name))
	osBody.SetAttributeValue("version", cty.StringVal(os.Version))
	osBody.SetAttributeValue("arch", cty.StringVal(os.Arch))
	if os.Kernel != "" {
		osBody.SetAttributeValue("kernel", cty.StringVal(os.Kernel))
	}
	if os.Platform != "" {
		osBody.SetAttributeValue("platform", cty.StringVal(os.Platform))
	}
	if os.Family != "" {
		osBody.SetAttributeValue("family", cty.StringVal(os.Family))
	}
}

// addHardwareFactsToHCL adds hardware facts to the HCL body
func (m *Manager) addHardwareFactsToHCL(systemBody *hclwrite.Body, hardware *spookytypesfacts.HardwareFacts) {
	hardwareBlock := systemBody.AppendNewBlock("hardware", nil)
	hardwareBody := hardwareBlock.Body()

	// CPU facts
	if hardware.CPU != nil {
		m.addCPUFactsToHCL(hardwareBody, hardware.CPU)
	}

	// Memory facts
	if hardware.Memory != nil {
		m.addMemoryFactsToHCL(hardwareBody, hardware.Memory)
	}
}

// addCPUFactsToHCL adds CPU facts to the HCL body
func (m *Manager) addCPUFactsToHCL(hardwareBody *hclwrite.Body, cpu *spookytypesfacts.CPUFacts) {
	cpuBlock := hardwareBody.AppendNewBlock("cpu", nil)
	cpuBody := cpuBlock.Body()
	cpuBody.SetAttributeValue("cores", cty.NumberIntVal(int64(cpu.Cores)))
	cpuBody.SetAttributeValue("model", cty.StringVal(cpu.Model))
	if cpu.Frequency > 0 {
		cpuBody.SetAttributeValue("frequency", cty.NumberFloatVal(cpu.Frequency))
	}
	if cpu.Architecture != "" {
		cpuBody.SetAttributeValue("architecture", cty.StringVal(cpu.Architecture))
	}
	if cpu.Vendor != "" {
		cpuBody.SetAttributeValue("vendor", cty.StringVal(cpu.Vendor))
	}
}

// addMemoryFactsToHCL adds memory facts to the HCL body
func (m *Manager) addMemoryFactsToHCL(hardwareBody *hclwrite.Body, memory *spookytypesfacts.MemoryFacts) {
	memoryBlock := hardwareBody.AppendNewBlock("memory", nil)
	memoryBody := memoryBlock.Body()
	memoryBody.SetAttributeValue("total", cty.NumberIntVal(memory.Total))
	if memory.Available > 0 {
		memoryBody.SetAttributeValue("available", cty.NumberIntVal(memory.Available))
	}
	if memory.Used > 0 {
		memoryBody.SetAttributeValue("used", cty.NumberIntVal(memory.Used))
	}
	if memory.Free > 0 {
		memoryBody.SetAttributeValue("free", cty.NumberIntVal(memory.Free))
	}
	if memory.Buffers > 0 {
		memoryBody.SetAttributeValue("buffers", cty.NumberIntVal(memory.Buffers))
	}
	if memory.Cached > 0 {
		memoryBody.SetAttributeValue("cached", cty.NumberIntVal(memory.Cached))
	}
}

// addNetworkFactsToHCL adds network facts to the HCL body
func (m *Manager) addNetworkFactsToHCL(systemBody *hclwrite.Body, network *spookytypesfacts.NetworkFacts) {
	networkBlock := systemBody.AppendNewBlock("network", nil)
	networkBody := networkBlock.Body()

	if network.Hostname != "" {
		networkBody.SetAttributeValue("hostname", cty.StringVal(network.Hostname))
	}

	if len(network.Interfaces) > 0 {
		m.addNetworkInterfacesToHCL(networkBody, network.Interfaces)
	}

	if len(network.IPAddresses) > 0 {
		// Convert IP addresses to HCL list
		ipList := make([]cty.Value, len(network.IPAddresses))
		for i, ip := range network.IPAddresses {
			ipList[i] = cty.StringVal(ip)
		}
		networkBody.SetAttributeValue("ip_addresses", cty.ListVal(ipList))
	}

	if network.PrimaryIP != "" {
		networkBody.SetAttributeValue("primary_ip", cty.StringVal(network.PrimaryIP))
	}
}

// addNetworkInterfacesToHCL adds network interfaces to the HCL body
func (m *Manager) addNetworkInterfacesToHCL(networkBody *hclwrite.Body, interfaces []*spookytypesfacts.NetworkInterface) {
	interfacesBlock := networkBody.AppendNewBlock("interfaces", nil)
	interfacesBody := interfacesBlock.Body()

	for _, iface := range interfaces {
		ifaceBlock := interfacesBody.AppendNewBlock("interface", []string{iface.Name})
		ifaceBody := ifaceBlock.Body()
		ifaceBody.SetAttributeValue("name", cty.StringVal(iface.Name))
		if iface.MACAddress != "" {
			ifaceBody.SetAttributeValue("mac_address", cty.StringVal(iface.MACAddress))
		}
		if len(iface.IPAddresses) > 0 {
			// Convert IP addresses to HCL list
			ipList := make([]cty.Value, len(iface.IPAddresses))
			for i, ip := range iface.IPAddresses {
				ipList[i] = cty.StringVal(ip)
			}
			ifaceBody.SetAttributeValue("ip_addresses", cty.ListVal(ipList))
		}
		if iface.MTU > 0 {
			ifaceBody.SetAttributeValue("mtu", cty.NumberIntVal(int64(iface.MTU)))
		}
	}
}

// CollectAndStoreFacts collects facts for a machine (in-memory only)
func (m *Manager) CollectAndStoreFacts(ctx context.Context, machine *spookytypes.Machine) error {
	if machine == nil {
		return fmt.Errorf("machine cannot be nil")
	}

	// Collect facts
	facts, err := m.CollectFacts(ctx, machine)
	if err != nil {
		return err
	}

	// Validate facts
	validationResult, err := m.ValidateFacts(ctx, facts)
	if err != nil {
		return fmt.Errorf("facts validation failed: %w", err)
	}

	// Type assert validation result
	if result, ok := validationResult.(*spookytypesschemas.ValidationResult); ok {
		if !result.Valid {
			return fmt.Errorf("facts validation failed: %v", result.Errors)
		}
	}

	// Facts are stored in memory for the duration of the operation
	m.logger.Info("Facts collected and stored in memory", map[string]interface{}{
		"machine": machine.Hostname,
	})

	return nil
}

// CollectAndStoreFactsParallel collects facts for multiple machines in parallel (in-memory only)
func (m *Manager) CollectAndStoreFactsParallel(ctx context.Context, machines []*spookytypes.Machine, maxWorkers int) error {
	if len(machines) == 0 {
		return fmt.Errorf("no machines provided")
	}

	if maxWorkers <= 0 {
		maxWorkers = 4
	}

	m.logger.Info("Starting parallel fact collection", map[string]interface{}{"machines": len(machines), "workers": maxWorkers})

	// Create worker pool
	semaphore := make(chan struct{}, maxWorkers)
	results := make(chan error, len(machines))

	// Start workers
	for _, machine := range machines {
		go func(machine *spookytypes.Machine) {
			semaphore <- struct{}{} // Acquire semaphore
			defer func() {
				<-semaphore // Release semaphore
			}()

			err := m.CollectAndStoreFacts(ctx, machine)
			results <- err
		}(machine)
	}

	// Collect results
	var errors []error
	for i := 0; i < len(machines); i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}

	// Report results
	successful := len(machines) - len(errors)
	m.logger.Info("Completed parallel fact collection", map[string]interface{}{"successful": successful, "failed": len(errors)})

	if len(errors) > 0 {
		return fmt.Errorf("fact collection failed for %d machines", len(errors))
	}

	return nil
}

// GetStorageStats returns memory storage statistics
func (m *Manager) GetStorageStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"storage_type": "in_memory_only",
		"description":  "Facts are stored in memory for the duration of operations only",
	}, nil
}

// GetFacts retrieves facts for a specific machine (collects on demand)
func (m *Manager) GetFacts(ctx context.Context, machineID string) (*spookytypes.FactCollection, error) {
	m.logger.Info("Getting facts for machine", map[string]interface{}{
		"machine": machineID,
	})

	// Facts are collected on demand since there's no persistent storage
	return m.CollectFacts(ctx, &spookytypes.Machine{Hostname: machineID})
}

// FactExport represents exported facts data
type FactExport struct {
	ExportedAt   time.Time                              `json:"exported_at"`
	Format       string                                 `json:"format"`
	MachineCount int                                    `json:"machine_count"`
	Facts        map[string]*spookytypes.FactCollection `json:"facts"`
}

// ClearFacts removes all facts from memory
func (m *Manager) ClearFacts(_ context.Context) error {
	m.logger.Info("Clearing all facts from memory")

	// In-memory storage - no cleanup needed as facts are temporary
	m.logger.Info("Successfully cleared all facts from memory")

	return nil
}

// DecryptFacts decrypts age-encrypted values in facts collection
func (m *Manager) DecryptFacts(ctx context.Context, facts *spookytypes.FactCollection, secretsIntegration spookyinterfaces.SecretsIntegration, identityPath string) error {
	if facts == nil {
		return fmt.Errorf("facts cannot be nil")
	}

	if secretsIntegration == nil {
		return fmt.Errorf("secrets integration cannot be nil")
	}

	if identityPath == "" {
		return fmt.Errorf("identity path cannot be empty")
	}

	m.logger.Info("Decrypting age-encrypted values in facts", map[string]interface{}{
		"machine_id": facts.MachineID,
	})

	// Use the HCL processor
	hclProcessor := spookysecrets.NewHCLProcessor(m.logger)
	err := hclProcessor.DecryptHCLValues(ctx, facts, secretsIntegration, identityPath)
	if err != nil {
		return fmt.Errorf("failed to decrypt facts: %w", err)
	}

	m.logger.Info("Successfully decrypted age-encrypted values in facts using HCL processor", map[string]interface{}{
		"machine_id": facts.MachineID,
	})

	return nil
}
