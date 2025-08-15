// Package types provides unified access to all spooky type definitions.
// This package re-exports types from domain-specific subpackages for convenient access.
package types

import (
	// Action types
	spookytypesactions "spooky/internal/types/actions"

	// CLI types
	spookytypescli "spooky/internal/types/cli"

	// Common types
	spookytypescommon "spooky/internal/types/common"

	// Configuration types
	spookytypesconfig "spooky/internal/types/config"

	// Facts types
	spookytypesfacts "spooky/internal/types/facts"

	// Logging types
	spookytypeslogging "spooky/internal/types/logging"

	// Machine types
	spookytypesmachines "spooky/internal/types/machines"

	// Project types
	spookytypesproject "spooky/internal/types/project"

	// Schema types
	spookytypesschemas "spooky/internal/types/schemas"

	// Template types
	spookytypestemplates "spooky/internal/types/templates"

	// Variable types
	spookytypesvariables "spooky/internal/types/variables"

	// SSH types
	spookyssh "spooky/internal/types/ssh"
)

// =============================================================================
// Action Types
// =============================================================================

// Action is a type alias for spookytypesactions.Action
type Action = spookytypesactions.Action

// ActionCollection is a type alias for spookytypesactions.ActionCollection
type ActionCollection = spookytypesactions.ActionCollection

// ActionPlan is a type alias for spookytypesactions.ActionPlan
type ActionPlan = spookytypesactions.ActionPlan

// ActionDependency is a type alias for spookytypesactions.ActionDependency
type ActionDependency = spookytypesactions.ActionDependency

// ActionRun is a type alias for spookytypesactions.ActionRun
type ActionRun = spookytypesactions.ActionRun

// ActionMetrics is a type alias for spookytypesactions.ActionMetrics
type ActionMetrics = spookytypesactions.ActionMetrics

// ActionValidation is a type alias for spookytypesactions.ActionValidation
type ActionValidation = spookytypesactions.ActionValidation

// ActionRunContext is a type alias for spookytypesactions.ActionRunContext
type ActionRunContext = spookytypesactions.ActionRunContext

// ActingSession is a type alias for spookytypesactions.ActingSession
type ActingSession = spookytypesactions.ActingSession

// ActingResult is a type alias for spookytypesactions.ActingResult
type ActingResult = spookytypesactions.ActingResult

// ActionError is a type alias for spookytypesactions.ActionError
type ActionError = spookytypesactions.ActionError

// ActingError is a type alias for spookytypesactions.ActingError
type ActingError = spookytypesactions.ActingError

// PlanningError is a type alias for spookytypesactions.PlanningError
type PlanningError = spookytypesactions.PlanningError

// ActionValidationError is a type alias for spookytypesactions.ValidationError
type ActionValidationError = spookytypesactions.ValidationError

// DependencyError is a type alias for spookytypesactions.DependencyError
type DependencyError = spookytypesactions.DependencyError

// ActionConfigurationError is a type alias for spookytypesactions.ConfigurationError
type ActionConfigurationError = spookytypesactions.ConfigurationError

// CommandConfig is a type alias for spookytypesactions.CommandConfig
type CommandConfig = spookytypesactions.CommandConfig

// ScriptConfig is a type alias for spookytypesactions.ScriptConfig
type ScriptConfig = spookytypesactions.ScriptConfig

// TemplateConfig is a type alias for spookytypesactions.TemplateConfig
type TemplateConfig = spookytypesactions.TemplateConfig

// FileCopyConfig is a type alias for spookytypesactions.FileCopyConfig
type FileCopyConfig = spookytypesactions.FileCopyConfig

// ServiceControlConfig is a type alias for spookytypesactions.ServiceControlConfig
type ServiceControlConfig = spookytypesactions.ServiceControlConfig

// ResourceLimits is a type alias for spookytypesactions.ResourceLimits
type ResourceLimits = spookytypesactions.ResourceLimits

// CollectionConfig is a type alias for spookytypesactions.CollectionConfig
type CollectionConfig = spookytypesactions.CollectionConfig

// ResourceUsage is a type alias for spookytypesactions.ResourceUsage
type ResourceUsage = spookytypesactions.ResourceUsage

// PerformancePoint is a type alias for spookytypesactions.PerformancePoint
type PerformancePoint = spookytypesactions.PerformancePoint

// ValidationWarning is a type alias for spookytypesactions.ValidationWarning
type ValidationWarning = spookytypesactions.ValidationWarning

// =============================================================================
// CLI Types
// =============================================================================

// Command is a type alias for spookytypescli.Command
type Command = spookytypescli.Command

// CommandContext is a type alias for spookytypescli.CommandContext
type CommandContext = spookytypescli.CommandContext

// CommandFlags is a type alias for spookytypescli.CommandFlags
type CommandFlags = spookytypescli.CommandFlags

// =============================================================================
// Common Types
// =============================================================================

// TimestampedEntity is a type alias for spookytypescommon.TimestampedEntity
type TimestampedEntity = spookytypescommon.TimestampedEntity

// NamedEntity is a type alias for spookytypescommon.NamedEntity
type NamedEntity = spookytypescommon.NamedEntity

// MetadataEntity is a type alias for spookytypescommon.MetadataEntity
type MetadataEntity = spookytypescommon.MetadataEntity

// ValidationEntity is a type alias for spookytypescommon.ValidationEntity
type ValidationEntity = spookytypescommon.ValidationEntity

// StatusEntity is a type alias for spookytypescommon.StatusEntity
type StatusEntity = spookytypescommon.StatusEntity

// CompleteEntity is a type alias for spookytypescommon.CompleteEntity
type CompleteEntity = spookytypescommon.CompleteEntity

// ErrorDetails is a type alias for spookytypescommon.ErrorDetails
type ErrorDetails = spookytypescommon.ErrorDetails

// ExportOptions is a type alias for spookytypescommon.ExportOptions
type ExportOptions = spookytypescommon.ExportOptions

// ImportOptions is a type alias for spookytypescommon.ImportOptions
type ImportOptions = spookytypescommon.ImportOptions

// Query is a type alias for spookytypescommon.Query
type Query = spookytypescommon.Query

// Result is a type alias for spookytypescommon.Result
type Result = spookytypescommon.Result

// TimeRange is a type alias for spookytypescommon.TimeRange
type TimeRange = spookytypescommon.TimeRange

// Pagination is a type alias for spookytypescommon.Pagination
type Pagination = spookytypescommon.Pagination

// EncryptionMetadata is a type alias for spookytypescommon.EncryptionMetadata
type EncryptionMetadata = spookytypescommon.EncryptionMetadata

// =============================================================================
// Configuration Types
// =============================================================================

// Config is a type alias for spookytypesconfig.Config
type Config = spookytypesconfig.Config

// =============================================================================
// Facts Types
// =============================================================================

// FactCollection is a type alias for spookytypesfacts.FactCollection
type FactCollection = spookytypesfacts.FactCollection

// =============================================================================
// Logging Types
// =============================================================================

// LogConfig is a type alias for spookytypeslogging.LogConfig
type LogConfig = spookytypeslogging.LogConfig

// LogEntry is a type alias for spookytypeslogging.LogEntry
type LogEntry = spookytypeslogging.LogEntry

// Logger is a type alias for spookytypeslogging.Logger
type Logger = spookytypeslogging.Logger

// LogLevel is a type alias for spookytypeslogging.LogLevel
type LogLevel = spookytypeslogging.LogLevel

// LogManager is a type alias for spookytypeslogging.LogManager
type LogManager = spookytypeslogging.LogManager

// =============================================================================
// Machine Types
// =============================================================================

// Machine is a type alias for spookytypesmachines.Machine
type Machine = spookytypesmachines.Machine

// MachineStatus is a type alias for spookytypesmachines.MachineStatus
type MachineStatus = spookytypesmachines.MachineStatus

// =============================================================================
// Project Types
// =============================================================================

// Project is a type alias for spookytypesproject.Project
type Project = spookytypesproject.Project

// ProjectMetadata is a type alias for spookytypesproject.Metadata
type ProjectMetadata = spookytypesproject.Metadata

// ProjectSettings is a type alias for spookytypesproject.Settings
type ProjectSettings = spookytypesproject.Settings

// ProjectConfig is a type alias for spookytypesproject.Config
type ProjectConfig = spookytypesproject.Config

// =============================================================================
// Schema Types
// =============================================================================

// Schema is a type alias for spookytypesschemas.Schema
type Schema = spookytypesschemas.Schema

// SchemaValidation is a type alias for spookytypesschemas.SchemaValidation
type SchemaValidation = spookytypesschemas.SchemaValidation

// SchemaError is a type alias for spookytypesschemas.SchemaError
type SchemaError = spookytypesschemas.SchemaError

// ValidationResult is a type alias for spookytypesschemas.ValidationResult
type ValidationResult = spookytypesschemas.ValidationResult

// FieldValidation is a type alias for spookytypesschemas.FieldValidation
type FieldValidation = spookytypesschemas.FieldValidation

// FieldConstraints is a type alias for spookytypesschemas.FieldConstraints
type FieldConstraints = spookytypesschemas.FieldConstraints

// CrossFieldValidation is a type alias for spookytypesschemas.CrossFieldValidation
type CrossFieldValidation = spookytypesschemas.CrossFieldValidation

// ValidationRule is a type alias for spookytypesschemas.ValidationRule
type ValidationRule = spookytypesschemas.ValidationRule

// ValidationErrorHandling is a type alias for spookytypesschemas.ValidationErrorHandling
type ValidationErrorHandling = spookytypesschemas.ValidationErrorHandling

// SchemaEvolution is a type alias for spookytypesschemas.SchemaEvolution
type SchemaEvolution = spookytypesschemas.SchemaEvolution

// SchemaVersion is a type alias for spookytypesschemas.SchemaVersion
type SchemaVersion = spookytypesschemas.SchemaVersion

// VersionChange is a type alias for spookytypesschemas.VersionChange
type VersionChange = spookytypesschemas.VersionChange

// VersionCompatibility is a type alias for spookytypesschemas.VersionCompatibility
type VersionCompatibility = spookytypesschemas.VersionCompatibility

// SchemaMigration is a type alias for spookytypesschemas.SchemaMigration
type SchemaMigration = spookytypesschemas.SchemaMigration

// MigrationValidation is a type alias for spookytypesschemas.MigrationValidation
type MigrationValidation = spookytypesschemas.MigrationValidation

// SchemaDeprecation is a type alias for spookytypesschemas.SchemaDeprecation
type SchemaDeprecation = spookytypesschemas.SchemaDeprecation

// BreakingChange is a type alias for spookytypesschemas.BreakingChange
type BreakingChange = spookytypesschemas.BreakingChange

// EvolutionPolicy is a type alias for spookytypesschemas.EvolutionPolicy
type EvolutionPolicy = spookytypesschemas.EvolutionPolicy

// SchemaCompatibility is a type alias for spookytypesschemas.SchemaCompatibility
type SchemaCompatibility = spookytypesschemas.SchemaCompatibility

// ErrorLocation is a type alias for spookytypesschemas.ErrorLocation
type ErrorLocation = spookytypesschemas.ErrorLocation

// CompatibilityResult is a type alias for spookytypesschemas.CompatibilityResult
type CompatibilityResult = spookytypesschemas.CompatibilityResult

// CompatibilityIssue is a type alias for spookytypesschemas.CompatibilityIssue
type CompatibilityIssue = spookytypesschemas.CompatibilityIssue

// ValidationStatistics is a type alias for spookytypesschemas.ValidationStatistics
type ValidationStatistics = spookytypesschemas.ValidationStatistics

// SchemaUpdate is a type alias for spookytypesschemas.SchemaUpdate
type SchemaUpdate = spookytypesschemas.SchemaUpdate

// SchemaAnalysis is a type alias for spookytypesschemas.SchemaAnalysis
type SchemaAnalysis = spookytypesschemas.SchemaAnalysis

// SchemaComplexity is a type alias for spookytypesschemas.SchemaComplexity
type SchemaComplexity = spookytypesschemas.SchemaComplexity

// SchemaCoverage is a type alias for spookytypesschemas.SchemaCoverage
type SchemaCoverage = spookytypesschemas.SchemaCoverage

// SchemaQuality is a type alias for spookytypesschemas.SchemaQuality
type SchemaQuality = spookytypesschemas.SchemaQuality

// SchemaComparison is a type alias for spookytypesschemas.SchemaComparison
type SchemaComparison = spookytypesschemas.SchemaComparison

// =============================================================================
// Template Types
// =============================================================================

// Template is a type alias for spookytypestemplates.Template
type Template = spookytypestemplates.Template

// TemplateContext is a type alias for spookytypestemplates.TemplateContext
type TemplateContext = spookytypestemplates.TemplateContext

// =============================================================================
// Variable Types
// =============================================================================

// Variable is a type alias for spookytypesvariables.Variable
type Variable = spookytypesvariables.Variable

// VariableContext is a type alias for spookytypesvariables.VariableContext
type VariableContext = spookytypesvariables.VariableContext

// VariableResolutionResult is a type alias for spookytypesvariables.VariableResolutionResult
type VariableResolutionResult = spookytypesvariables.VariableResolutionResult

// VariableValidation is a type alias for spookytypesvariables.VariableValidation
type VariableValidation = spookytypesvariables.VariableValidation

// VariableDependency is a type alias for spookytypesvariables.VariableDependency
type VariableDependency = spookytypesvariables.VariableDependency

// =============================================================================
// SSH Types
// =============================================================================

// SSHError is a type alias for spookyssh.Error
type SSHError = spookyssh.Error

// SSHErrorType is a type alias for spookyssh.ErrorType
type SSHErrorType = spookyssh.ErrorType

// ConnectionError is a type alias for spookyssh.ConnectionError
type ConnectionError = spookyssh.ConnectionError

// AuthenticationError is a type alias for spookyssh.AuthenticationError
type AuthenticationError = spookyssh.AuthenticationError

// HostKeyError is a type alias for spookyssh.HostKeyError
type HostKeyError = spookyssh.HostKeyError

// CommandError is a type alias for spookyssh.CommandError
type CommandError = spookyssh.CommandError

// SessionError is a type alias for spookyssh.SessionError
type SessionError = spookyssh.SessionError

// FileTransferError is a type alias for spookyssh.FileTransferError
type FileTransferError = spookyssh.FileTransferError

// ValidationError is a type alias for spookyssh.ValidationError
type ValidationError = spookyssh.ValidationError

// ConfigurationError is a type alias for spookyssh.ConfigurationError
type ConfigurationError = spookyssh.ConfigurationError

// ClientConfig is a type alias for spookyssh.ClientConfig
type ClientConfig = spookyssh.ClientConfig

// Client is a type alias for spookyssh.Client
type Client = spookyssh.Client

// ConnectionRequest is a type alias for spookyssh.ConnectionRequest
type ConnectionRequest = spookyssh.ConnectionRequest

// ConnectionResult is a type alias for spookyssh.ConnectionResult
type ConnectionResult = spookyssh.ConnectionResult

// Connection is a type alias for spookyssh.Connection
type Connection = spookyssh.Connection

// Authentication is a type alias for spookyssh.Authentication
type Authentication = spookyssh.Authentication

// AuthenticationResult is a type alias for spookyssh.AuthenticationResult
type AuthenticationResult = spookyssh.AuthenticationResult

// Session is a type alias for spookyssh.Session
type Session = spookyssh.Session

// Command is a type alias for spookyssh.Command
type SSHCommand = spookyssh.Command

// CommandResult is a type alias for spookyssh.CommandResult
type SSHCommandResult = spookyssh.CommandResult

// FileTransfer is a type alias for spookyssh.FileTransfer
type FileTransfer = spookyssh.FileTransfer

// FileTransferResult is a type alias for spookyssh.FileTransferResult
type FileTransferResult = spookyssh.FileTransferResult

// ConnectionPool is a type alias for spookyssh.ConnectionPool
type ConnectionPool = spookyssh.ConnectionPool

// AuthMethod is a type alias for spookyssh.AuthMethod
type AuthMethod = spookyssh.AuthMethod

// KeyType is a type alias for spookyssh.KeyType
type KeyType = spookyssh.KeyType

// TransferMode is a type alias for spookyssh.TransferMode
type TransferMode = spookyssh.TransferMode

// TransferDirection is a type alias for spookyssh.TransferDirection
type TransferDirection = spookyssh.TransferDirection

// PtyConfig is a type alias for spookyssh.PtyConfig
type PtyConfig = spookyssh.PtyConfig

// ClientStatus is a type alias for spookyssh.ClientStatus
type ClientStatus = spookyssh.ClientStatus

// SessionStatus is a type alias for spookyssh.SessionStatus
type SessionStatus = spookyssh.SessionStatus

// ConnectionStatus is a type alias for spookyssh.ConnectionStatus
type ConnectionStatus = spookyssh.ConnectionStatus

// TrustLevel is a type alias for spookyssh.TrustLevel
type TrustLevel = spookyssh.TrustLevel

// ValidationType is a type alias for spookyssh.ValidationType
type ValidationType = spookyssh.ValidationType

// ValidationLevel is a type alias for spookyssh.ValidationLevel
type ValidationLevel = spookyssh.ValidationLevel

// ActingSessionStatus is a type alias for spookyssh.ActingSessionStatus
type ActingSessionStatus = spookyssh.ActingSessionStatus

// ActingMode is a type alias for spookyssh.ActingMode
type ActingMode = spookyssh.ActingMode

// ActingBatchStatus is a type alias for spookyssh.ActingBatchStatus
type ActingBatchStatus = spookyssh.ActingBatchStatus

// ActingCommand is a type alias for spookyssh.ActingCommand
type ActingCommand = spookyssh.ActingCommand

// ActingCommandResult is a type alias for spookyssh.ActingCommandResult
type ActingCommandResult = spookyssh.ActingCommandResult

// ActingBatch is a type alias for spookyssh.ActingBatch
type ActingBatch = spookyssh.ActingBatch

// ActingBatchResult is a type alias for spookyssh.ActingBatchResult
type ActingBatchResult = spookyssh.ActingBatchResult

// ActingScript is a type alias for spookyssh.ActingScript
type ActingScript = spookyssh.ActingScript

// ActingScriptResult is a type alias for spookyssh.ActingScriptResult
type ActingScriptResult = spookyssh.ActingScriptResult

// Key is a type alias for spookyssh.Key
type Key = spookyssh.Key

// KeyPair is a type alias for spookyssh.KeyPair
type KeyPair = spookyssh.KeyPair

// HostKey is a type alias for spookyssh.HostKey
type HostKey = spookyssh.HostKey

// KnownHosts is a type alias for spookyssh.KnownHosts
type KnownHosts = spookyssh.KnownHosts

// ConnectionFactory is a type alias for spookyssh.ConnectionFactory
type ConnectionFactory = spookyssh.ConnectionFactory
