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

// ProjectMetadata is a type alias for spookytypesproject.ProjectMetadata
type ProjectMetadata = spookytypesproject.ProjectMetadata

// ProjectSettings is a type alias for spookytypesproject.ProjectSettings
type ProjectSettings = spookytypesproject.ProjectSettings

// ProjectConfig is a type alias for spookytypesproject.ProjectConfig
type ProjectConfig = spookytypesproject.ProjectConfig

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
