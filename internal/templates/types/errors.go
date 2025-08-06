package types

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `hcl:"field" json:"field"`
	Message string `hcl:"message" json:"message"`
}

// TemplateError represents a template error
type TemplateError struct {
	Template string `hcl:"template" json:"template"`
	Error    string `hcl:"error" json:"error"`
}

// FunctionError represents a function error
type FunctionError struct {
	Function string `hcl:"function" json:"function"`
	Error    string `hcl:"error" json:"error"`
}
