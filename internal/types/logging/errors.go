package logging

import "errors"

// Logging-related errors
var (
	ErrInvalidConfig = errors.New("invalid logging configuration")
	ErrBackendFailed = errors.New("logging backend failed to initialize")
	ErrOutputFailed  = errors.New("logging output failed to initialize")
)
