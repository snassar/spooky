package facts

// OpenTofuCollector collects facts from OpenTofu state files and outputs
type OpenTofuCollector struct {
	*StubCollector
}

// NewOpenTofuCollector creates a new OpenTofu fact collector
func NewOpenTofuCollector() *OpenTofuCollector {
	return &OpenTofuCollector{
		StubCollector: NewStubCollector("OpenTofu"),
	}
}
