// Package types provides unified access to all type definitions in the spooky codebase.
// This package serves as the central repository for all type definitions and re-exports
// types from domain-specific subpackages for consistent access patterns.
package types

import (
	spookytypescli "spooky/internal/types/cli"
	spookytypescommon "spooky/internal/types/common"
	spookytypesconfig "spooky/internal/types/config"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesproject "spooky/internal/types/project"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypesvariables "spooky/internal/types/variables"
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

// Logging types
type LogLevel = spookytypeslogging.LogLevel
type LogConfig = spookytypeslogging.LogConfig
type LogRotationConfig = spookytypeslogging.LogRotationConfig
type LogFilteringConfig = spookytypeslogging.LogFilteringConfig
type LogPerformanceConfig = spookytypeslogging.LogPerformanceConfig
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

// Machine types
type Machine = spookytypesmachines.Machine
type MachineStatus = spookytypesmachines.MachineStatus
type MachineCollection = spookytypesmachines.MachineCollection
type MachineFilter = spookytypesmachines.MachineFilter
type MachineQuery = spookytypesmachines.MachineQuery
type MachineResult = spookytypesmachines.MachineResult

// Variable types
type Variable = spookytypesvariables.Variable
type VariableType = spookytypesvariables.VariableType
type VariableScope = spookytypesvariables.VariableScope
type VariableValidation = spookytypesvariables.VariableValidation
type VariableConstraints = spookytypesvariables.VariableConstraints
type VariableCollection = spookytypesvariables.VariableCollection
type VariableResolutionConfig = spookytypesvariables.VariableResolutionConfig
type VariableSecurityConfig = spookytypesvariables.VariableSecurityConfig
type VariableContext = spookytypesvariables.VariableContext
type VariableResolutionResult = spookytypesvariables.VariableResolutionResult
type VariableError = spookytypesvariables.VariableError
type VariableErrorType = spookytypesvariables.VariableErrorType
type VariableWarning = spookytypesvariables.VariableWarning
type VariableWarningType = spookytypesvariables.VariableWarningType
type VariableValidationResult = spookytypesvariables.VariableValidationResult
type VariableExportOptions = spookytypesvariables.VariableExportOptions
type VariableImportOptions = spookytypesvariables.VariableImportOptions
type VariableQuery = spookytypesvariables.VariableQuery
type VariableResult = spookytypesvariables.VariableResult
type VariableFile = spookytypesvariables.VariableFile
type VariableSource = spookytypesvariables.VariableSource
type VariableDependency = spookytypesvariables.VariableDependency
type VariableDependencyGraph = spookytypesvariables.VariableDependencyGraph
type VariableStatistics = spookytypesvariables.VariableStatistics

// Variable type constants
const (
	VariableTypeString   = spookytypesvariables.VariableTypeString
	VariableTypeNumber   = spookytypesvariables.VariableTypeNumber
	VariableTypeFloat    = spookytypesvariables.VariableTypeFloat
	VariableTypeBool     = spookytypesvariables.VariableTypeBool
	VariableTypeList     = spookytypesvariables.VariableTypeList
	VariableTypeMap      = spookytypesvariables.VariableTypeMap
	VariableTypeObject   = spookytypesvariables.VariableTypeObject
	VariableTypeDuration = spookytypesvariables.VariableTypeDuration
	VariableTypeIP       = spookytypesvariables.VariableTypeIP
	VariableTypeCIDR     = spookytypesvariables.VariableTypeCIDR
	VariableTypePath     = spookytypesvariables.VariableTypePath
	VariableTypeFile     = spookytypesvariables.VariableTypeFile
	VariableTypeSecret   = spookytypesvariables.VariableTypeSecret
)

// Variable scope constants
const (
	VariableScopeProject   = spookytypesvariables.VariableScopeProject
	VariableScopeGlobal    = spookytypesvariables.VariableScopeGlobal
	VariableScopeInherited = spookytypesvariables.VariableScopeInherited
)

// Variable error type constants
const (
	VariableErrorTypeValidation   = spookytypesvariables.VariableErrorTypeValidation
	VariableErrorTypeResolution   = spookytypesvariables.VariableErrorTypeResolution
	VariableErrorTypeDependency   = spookytypesvariables.VariableErrorTypeDependency
	VariableErrorTypeCircular     = spookytypesvariables.VariableErrorTypeCircular
	VariableErrorTypeMissing      = spookytypesvariables.VariableErrorTypeMissing
	VariableErrorTypeTypeMismatch = spookytypesvariables.VariableErrorTypeTypeMismatch
	VariableErrorTypeConstraint   = spookytypesvariables.VariableErrorTypeConstraint
	VariableErrorTypeSecurity     = spookytypesvariables.VariableErrorTypeSecurity
)

// Variable warning type constants
const (
	VariableWarningTypeDeprecated   = spookytypesvariables.VariableWarningTypeDeprecated
	VariableWarningTypeUnused       = spookytypesvariables.VariableWarningTypeUnused
	VariableWarningTypeSensitive    = spookytypesvariables.VariableWarningTypeSensitive
	VariableWarningTypeUnencrypted  = spookytypesvariables.VariableWarningTypeUnencrypted
	VariableWarningTypeDefaultValue = spookytypesvariables.VariableWarningTypeDefaultValue
	VariableWarningTypeEnvironment  = spookytypesvariables.VariableWarningTypeEnvironment
	VariableWarningTypeDependency   = spookytypesvariables.VariableWarningTypeDependency
)

// Placeholder types for interfaces (to be implemented)
type FactCollection = interface{}
type FactStorage = interface{}
type Action = interface{}
type ActingResult = interface{}
type Template = interface{}

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
