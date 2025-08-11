// Package types provides unified access to all type definitions in the spooky codebase.
// This package serves as the central repository for all type definitions and re-exports
// types from domain-specific subpackages for consistent access patterns.
package types

import (
	spookytypescli "spooky/internal/types/cli"
	spookytypescommon "spooky/internal/types/common"
	spookytypesconfig "spooky/internal/types/config"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesproject "spooky/internal/types/project"
	spookytypesschemas "spooky/internal/types/schemas"
)

// Re-export all types from subpackages for unified access
// This enables consistent import patterns across the codebase

// CLI types
type Command = spookytypescli.Command
type CommandContext = spookytypescli.CommandContext
type CommandFlags = spookytypescli.CommandFlags
type CommandError = spookytypescli.CommandError

// Common types
type TimestampedEntity = spookytypescommon.TimestampedEntity
type NamedEntity = spookytypescommon.NamedEntity
type MetadataEntity = spookytypescommon.MetadataEntity
type ValidationEntity = spookytypescommon.ValidationEntity
type StatusEntity = spookytypescommon.StatusEntity
type CompleteEntity = spookytypescommon.CompleteEntity
type ErrorDetails = spookytypescommon.ErrorDetails
type ExportOptions = spookytypescommon.ExportOptions
type ImportOptions = spookytypescommon.ImportOptions
type Query = spookytypescommon.Query
type Result = spookytypescommon.Result
type TimeRange = spookytypescommon.TimeRange
type Pagination = spookytypescommon.Pagination
type EncryptionMetadata = spookytypescommon.EncryptionMetadata

// Config types
type Config = spookytypesconfig.Config
type GlobalConfig = spookytypesconfig.GlobalConfig
type CLIConfig = spookytypesconfig.CLIConfig
type LoggingConfig = spookytypesconfig.LoggingConfig
type SSHConfig = spookytypesconfig.SSHConfig
type StorageConfig = spookytypesconfig.StorageConfig
type SecurityConfig = spookytypesconfig.SecurityConfig
type PerformanceConfig = spookytypesconfig.PerformanceConfig
type IsolationConfig = spookytypesconfig.IsolationConfig

// Logging types
type LogLevel = spookytypeslogging.LogLevel
type LogConfig = spookytypeslogging.LogConfig
type LogRotation = spookytypeslogging.LogRotation
type LogFiltering = spookytypeslogging.LogFiltering
type LogPerformance = spookytypeslogging.LogPerformance
type LogEntry = spookytypeslogging.LogEntry
type LogError = spookytypeslogging.LogError
type LogCaller = spookytypeslogging.LogCaller
type Logger = spookytypeslogging.Logger
type LogManager = spookytypeslogging.LogManager

// Log level constants
const (
	LogLevelDebug = spookytypeslogging.LogLevelDebug
	LogLevelInfo  = spookytypeslogging.LogLevelInfo
	LogLevelWarn  = spookytypeslogging.LogLevelWarn
	LogLevelError = spookytypeslogging.LogLevelError
	LogLevelFatal = spookytypeslogging.LogLevelFatal
)

// Project types
type Project = spookytypesproject.Project
type ProjectConfig = spookytypesproject.ProjectConfig
type ProjectMetadata = spookytypesproject.ProjectMetadata
type ProjectSettings = spookytypesproject.ProjectSettings
type ProjectValidation = spookytypesproject.ProjectValidation
type ProjectStatus = spookytypesproject.ProjectStatus
type ProjectIsolation = spookytypesproject.ProjectIsolation
type ProjectExecution = spookytypesproject.ProjectExecution
type ProjectManager = spookytypesproject.ProjectManager
type ProjectValidator = spookytypesproject.ProjectValidator
type ProjectLoader = spookytypesproject.ProjectLoader

// Schema types
type Schema = spookytypesschemas.Schema
type SchemaValidation = spookytypesschemas.SchemaValidation
type ValidationRule = spookytypesschemas.ValidationRule
type ValidationErrorHandling = spookytypesschemas.ValidationErrorHandling
type SchemaError = spookytypesschemas.SchemaError
type SchemaRegistry = spookytypesschemas.SchemaRegistry
type ValidationResult = spookytypesschemas.ValidationResult
type SchemaValidator = spookytypesschemas.SchemaValidator
type SchemaLoader = spookytypesschemas.SchemaLoader

// Schema error constructor
var NewSchemaError = spookytypesschemas.NewSchemaError
