package facts

// HCLCollector collects facts from HCL configuration files
type HCLCollector struct {
	*StubCollector
}

// NewHCLCollector creates a new HCL fact collector
func NewHCLCollector() *HCLCollector {
	return &HCLCollector{
		StubCollector: NewStubCollector("HCL"),
	}
}
