// Package types provides all type definitions for the spooky project.
// This package consolidates all types from various domains into a unified structure.
package types

// Re-export all types from subpackages for unified access
import (
	// Import subpackages to access their types
	spookytypesactions "spooky/internal/types/actions"
	spookytypescli "spooky/internal/types/cli"
	spookytypescommon "spooky/internal/types/common"
	spookytypesconfig "spooky/internal/types/config"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesproject "spooky/internal/types/project"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypessecrets "spooky/internal/types/secrets"
	spookytypesssh "spooky/internal/types/ssh"
	spookytypestemplates "spooky/internal/types/templates"
	spookytypesvariables "spooky/internal/types/variables"
)

// Re-export facts types
type Fact = spookytypesfacts.Fact
type FactCollection = spookytypesfacts.FactCollection
type FactSource = spookytypesfacts.FactSource
type MergePolicy = spookytypesfacts.MergePolicy
type MergeMode = spookytypesfacts.MergeMode
type CustomFacts = spookytypesfacts.CustomFacts
type ImportOptions = spookytypesfacts.ImportOptions
type SystemFacts = spookytypesfacts.SystemFacts
type OSInfo = spookytypesfacts.OSInfo
type HardwareInfo = spookytypesfacts.HardwareInfo
type CPUInfo = spookytypesfacts.CPUInfo
type MemoryInfo = spookytypesfacts.MemoryInfo
type StorageInfo = spookytypesfacts.StorageInfo
type DiskInfo = spookytypesfacts.DiskInfo
type NetworkInfo = spookytypesfacts.NetworkInfo
type InterfaceInfo = spookytypesfacts.InterfaceInfo
type DNSInfo = spookytypesfacts.DNSInfo
type ExportOptions = spookytypesfacts.ExportOptions
type TimeRange = spookytypesfacts.TimeRange
type FactQuery = spookytypesfacts.FactQuery
type StorageType = spookytypesfacts.StorageType
type StorageOptions = spookytypesfacts.StorageOptions
type EncryptionMetadata = spookytypesfacts.EncryptionMetadata
type FactCollector = spookytypesfacts.FactCollector

// Re-export facts constants
const (
	SourceSSH      = spookytypesfacts.SourceSSH
	SourceLocal    = spookytypesfacts.SourceLocal
	SourceHCL      = spookytypesfacts.SourceHCL
	SourceOpenTofu = spookytypesfacts.SourceOpenTofu
	SourceCustom   = spookytypesfacts.SourceCustom
	SourceJSON     = spookytypesfacts.SourceJSON

	// Merge policy constants
	MergePolicyReplace = spookytypesfacts.MergePolicyReplace
	MergePolicyMerge   = spookytypesfacts.MergePolicyMerge
	MergePolicyAppend  = spookytypesfacts.MergePolicyAppend
	MergePolicySkip    = spookytypesfacts.MergePolicySkip

	// Merge mode constants
	MergeModeReplace = spookytypesfacts.MergeModeReplace
	MergeModeMerge   = spookytypesfacts.MergeModeMerge
	MergeModeAppend  = spookytypesfacts.MergeModeAppend
	MergeModeSelect  = spookytypesfacts.MergeModeSelect

	// Default TTL
	DefaultTTL = spookytypesfacts.DefaultTTL
)

// Re-export config types
type Config = spookytypesconfig.Config
type GlobalConfig = spookytypesconfig.GlobalConfig
type ProjectConfig = spookytypesconfig.ProjectConfig
type EnvironmentConfig = spookytypesconfig.EnvironmentConfig
type LoadingConfig = spookytypesconfig.LoadingConfig
type ConfigValidationConfig = spookytypesconfig.ValidationConfig
type ValidationRules = spookytypesconfig.ValidationRules
type ConfigSource = spookytypesconfig.ConfigSource
type ConfigValidationError = spookytypesconfig.ValidationError
type ConfigError = spookytypesconfig.ConfigError
type EnvironmentVariable = spookytypesconfig.EnvironmentVariable
type EnvironmentSettings = spookytypesconfig.EnvironmentSettings
type EnvironmentSource = spookytypesconfig.EnvironmentSource
type EnvironmentValue = spookytypesconfig.EnvironmentValue
type EnvironmentValidationRule = spookytypesconfig.EnvironmentValidationRule
type EnvironmentValidationResult = spookytypesconfig.EnvironmentValidationResult
type StorageConfig = spookytypesconfig.StorageConfig
type FactsConfig = spookytypesconfig.FactsConfig
type ConfigSSHConfig = spookytypesconfig.SSHConfig
type TemplatesConfig = spookytypesconfig.TemplatesConfig
type SecurityConfig = spookytypesconfig.SecurityConfig
type AgeConfig = spookytypesconfig.AgeConfig
type ConfigLoggingConfig = spookytypesconfig.LoggingConfig
type PerformanceConfig = spookytypesconfig.PerformanceConfig
type IsolationConfig = spookytypesconfig.IsolationConfig
type ConfigMachine = spookytypesconfig.Machine
type ConfigAction = spookytypesconfig.Action
type ConfigActionResourceLimits = spookytypesconfig.ActionResourceLimits
type ConfigTemplateConfig = spookytypesconfig.TemplateConfig

// Re-export actions types
type Action = spookytypesactions.Action
type ActionState = spookytypesactions.ActionState
type ActionStatus = spookytypesactions.ActionStatus
type ActionCollection = spookytypesactions.ActionCollection
type ActionTemplateConfig = spookytypesactions.TemplateConfig
type ActionResourceLimits = spookytypesactions.ActionResourceLimits
type ActionContext = spookytypesactions.ActionContext
type ActingContext = spookytypesactions.ActingContext
type ActingState = spookytypesactions.ActingState
type RunStatus = spookytypesactions.RunStatus
type ActingSession = spookytypesactions.ActingSession
type RunResult = spookytypesactions.RunResult
type CollectionPlan = spookytypesactions.CollectionPlan
type FactRequirements = spookytypesactions.FactRequirements
type FactValidationRule = spookytypesactions.FactValidationRule
type PlannedCollection = spookytypesactions.PlannedCollection
type PlannedAction = spookytypesactions.PlannedAction
type CollectionDependency = spookytypesactions.CollectionDependency
type DependencyType = spookytypesactions.DependencyType
type RollbackPlan = spookytypesactions.RollbackPlan
type CollectionValidation = spookytypesactions.CollectionValidation
type FactValidationResult = spookytypesactions.FactValidationResult
type RetryPolicy = spookytypesactions.RetryPolicy
type PlanningOptions = spookytypesactions.PlanningOptions
type ActionFactSource = spookytypesactions.FactSource
type ActionFactDependencies = spookytypesactions.ActionFactDependencies
type DependencyNode = spookytypesactions.DependencyNode
type OrchestrationPlan = spookytypesactions.OrchestrationPlan
type OrchestrationResult = spookytypesactions.OrchestrationResult
type OrchestrationStatus = spookytypesactions.OrchestrationStatus
type RollbackStrategy = spookytypesactions.RollbackStrategy
type ActionPlan = spookytypesactions.ActionPlan
type PerformanceMetrics = spookytypesactions.PerformanceMetrics

// Re-export machines types
type MachineSecrets = spookytypesmachines.MachineSecrets
type MachineKeys = spookytypesmachines.MachineKeys
type MachineSecretsConfig = spookytypesmachines.SecretsConfig
type MachineSecretsStatus = spookytypesmachines.SecretsStatus
type MachineKeyMetadata = spookytypesmachines.KeyMetadata
type MachineSecretsError = spookytypesmachines.MachineSecretsError
type ExportConfig = spookytypesmachines.ExportConfig
type ImportConfig = spookytypesmachines.ImportConfig
type ConflictResolution = spookytypesmachines.ConflictResolution
type ImportResult = spookytypesmachines.ImportResult
type ImportConflict = spookytypesmachines.ImportConflict
type ImportError = spookytypesmachines.ImportError
type MachineValidationError = spookytypesmachines.ValidationError
type MachineError = spookytypesmachines.MachineError
type IndexError = spookytypesmachines.IndexError
type ConnectivityError = spookytypesmachines.ConnectivityError
type ExportError = spookytypesmachines.ExportError
type IndexType = spookytypesmachines.IndexType
type ConnectivityTestPhase = spookytypesmachines.ConnectivityTestPhase
type ExportFormat = spookytypesmachines.ExportFormat
type IndexManagerState = spookytypesmachines.IndexManagerState
type IndexSyncData = spookytypesmachines.IndexSyncData
type IndexMetrics = spookytypesmachines.IndexMetrics
type IndexTypeStats = spookytypesmachines.IndexTypeStats
type IndexPerformanceStats = spookytypesmachines.IndexPerformanceStats
type IndexTypePerformanceStats = spookytypesmachines.IndexTypePerformanceStats
type ConnectivityTestResult = spookytypesmachines.ConnectivityTestResult
type ConnectivityTestSummary = spookytypesmachines.ConnectivityTestSummary
type PhaseSummary = spookytypesmachines.PhaseSummary
type MachineFilter = spookytypesmachines.MachineFilter
type ExportResult = spookytypesmachines.ExportResult
type ExportStats = spookytypesmachines.ExportStats
type EncryptionResult = spookytypesmachines.EncryptionResult
type EncryptionError = spookytypesmachines.EncryptionError
type MachineEncryptionStatus = spookytypesmachines.MachineEncryptionStatus
type FieldEncryptionStatus = spookytypesmachines.FieldEncryptionStatus
type MachineIndex = spookytypesmachines.MachineIndex
type IndexManagerConfig = spookytypesmachines.IndexManagerConfig
type ConnectivityTestOptions = spookytypesmachines.ConnectivityTestOptions
type MachineExportOptions = spookytypesmachines.ExportOptions
type MachineValidationConfig = spookytypesmachines.ValidationConfig
type MachineValidationResult = spookytypesmachines.ValidationResult
type ValidationStatus = spookytypesmachines.ValidationStatus
type ValidationHistory = spookytypesmachines.ValidationHistory
type ValidationSummary = spookytypesmachines.ValidationSummary

// Re-export machines constants
const (
	IndexTypeName     = spookytypesmachines.IndexTypeName
	IndexTypeHost     = spookytypesmachines.IndexTypeHost
	IndexTypeTag      = spookytypesmachines.IndexTypeTag
	IndexTypeGroup    = spookytypesmachines.IndexTypeGroup
	IndexTypeNetwork  = spookytypesmachines.IndexTypeNetwork
	IndexTypeUser     = spookytypesmachines.IndexTypeUser
	IndexTypePort     = spookytypesmachines.IndexTypePort
	IndexTypeMetadata = spookytypesmachines.IndexTypeMetadata

	PhaseDNS = spookytypesmachines.PhaseDNS
	PhaseSSH = spookytypesmachines.PhaseSSH

	ExportFormatHCL  = spookytypesmachines.ExportFormatHCL
	ExportFormatJSON = spookytypesmachines.ExportFormatJSON
)

// Re-export ssh types
type SSHConfig = spookytypesssh.SSHConfig
type SSHConnection = spookytypesssh.SSHConnection
type CommandResult = spookytypesssh.CommandResult
type ConnectionInfo = spookytypesssh.ConnectionInfo
type ConnectionState = spookytypesssh.ConnectionState
type HostKeyCallback = spookytypesssh.HostKeyCallback
type ClientConfig = spookytypesssh.ClientConfig
type SSHClientConfig = spookytypesssh.Config
type ConnectionStats = spookytypesssh.ConnectionStats
type SSHMachine = spookytypesssh.Machine
type SSHKey = spookytypesssh.SSHKey
type KeysConfig = spookytypesssh.KeysConfig
type KeyValidationResult = spookytypesssh.KeyValidationResult
type SSHAction = spookytypesssh.SSHAction
type ActionResult = spookytypesssh.ActionResult
type TemplateAction = spookytypesssh.TemplateAction
type ExecuteConfig = spookytypesssh.ExecuteConfig
type ActingConfig = spookytypesssh.ActingConfig
type SSHError = spookytypesssh.SSHError
type ConnectionError = spookytypesssh.ConnectionError
type AuthenticationError = spookytypesssh.AuthenticationError
type ExecutionError = spookytypesssh.ExecutionError
type PooledConnection = spookytypesssh.PooledConnection
type PoolConfig = spookytypesssh.PoolConfig
type PoolStats = spookytypesssh.PoolStats
type ConnectionEvent = spookytypesssh.ConnectionEvent
type ConnectionMetrics = spookytypesssh.ConnectionMetrics
type AuthenticationMethod = spookytypesssh.AuthenticationMethod
type AuthenticationConfig = spookytypesssh.AuthenticationConfig
type AuthenticationResult = spookytypesssh.AuthenticationResult
type AuthenticationCache = spookytypesssh.AuthenticationCache

// Re-export ssh constants
const (
	AuthMethodPassword = spookytypesssh.AuthMethodPassword
	AuthMethodKey      = spookytypesssh.AuthMethodKey
	AuthMethodMixed    = spookytypesssh.AuthMethodMixed
)

// Re-export templates types
type Template = spookytypestemplates.Template
type TemplateContext = spookytypestemplates.TemplateContext
type TemplateData = spookytypestemplates.TemplateData
type TemplateVariable = spookytypestemplates.TemplateVariable
type TemplateMetadata = spookytypestemplates.TemplateMetadata
type TemplateValidationError = spookytypestemplates.ValidationError
type TemplateError = spookytypestemplates.TemplateError
type FunctionError = spookytypestemplates.FunctionError
type FunctionsConfig = spookytypestemplates.FunctionsConfig
type EngineConfig = spookytypestemplates.EngineConfig
type TemplateValidationConfig = spookytypestemplates.ValidationConfig
type TemplateSecretsConfig = spookytypestemplates.SecretsConfig
type TemplateConfig = spookytypestemplates.Config
type TemplateValidationRules = spookytypestemplates.ValidationRules

// Re-export variables types
type Variable = spookytypesvariables.Variable
type VariableCollection = spookytypesvariables.VariableCollection
type VariableContext = spookytypesvariables.VariableContext
type ResolutionConfig = spookytypesvariables.ResolutionConfig
type DependencyGraph = spookytypesvariables.DependencyGraph
type VariableDependencyNode = spookytypesvariables.DependencyNode
type VariableExportFormat = spookytypesvariables.ExportFormat
type VariableImportFormat = spookytypesvariables.ImportFormat
type VariableExportOptions = spookytypesvariables.ExportOptions
type VariableImportOptions = spookytypesvariables.ImportOptions
type VariableValidationResult = spookytypesvariables.ValidationResult
type VariableValidationError = spookytypesvariables.ValidationError
type VariableValidationConfig = spookytypesvariables.ValidationConfig
type VariableValidationWarning = spookytypesvariables.ValidationWarning
type VariableValidationRules = spookytypesvariables.ValidationRules

// Re-export logging types
type LoggingConfig = spookytypeslogging.Config
type LogLevel = spookytypeslogging.LogLevel
type LogEntry = spookytypeslogging.LogEntry
type Field = spookytypeslogging.Field

// Re-export logging constants
const (
	DebugLevel = spookytypeslogging.DebugLevel
	InfoLevel  = spookytypeslogging.InfoLevel
	WarnLevel  = spookytypeslogging.WarnLevel
	ErrorLevel = spookytypeslogging.ErrorLevel
)

// Re-export secrets types
type EncryptedValue = spookytypessecrets.EncryptedValue
type SecretsError = spookytypessecrets.SecretsError
type SecretsConfig = spookytypessecrets.SecretsConfig
type SecretsKeyMetadata = spookytypessecrets.KeyMetadata
type SecretsStatus = spookytypessecrets.SecretsStatus

// Re-export schemas types
type SchemaType = spookytypesschemas.SchemaType
type Schema = spookytypesschemas.Schema
type SchemaConfig = spookytypesschemas.Config
type SchemaValidationError = spookytypesschemas.ValidationError
type SchemaValidationResult = spookytypesschemas.ValidationResult

// Re-export cli types
type Command = spookytypescli.Command
type Context = spookytypescli.Context
type ExecutionContext = spookytypescli.ExecutionContext
type CLIError = spookytypescli.CLIError
type Flag = spookytypescli.Flag
type FlagSet = spookytypescli.FlagSet
type FlagValue = spookytypescli.FlagValue

// Re-export project types
type Project = spookytypesproject.Project

// Re-export common types
type TimestampedEntity = spookytypescommon.TimestampedEntity
type NamedEntity = spookytypescommon.NamedEntity
type MetadataEntity = spookytypescommon.MetadataEntity
type ValidationEntity = spookytypescommon.ValidationEntity
type StatusEntity = spookytypescommon.StatusEntity
type CompleteEntity = spookytypescommon.CompleteEntity

// This file serves as the main entry point for all spookytypes.
// Import this package to access all types:
//
//
// Then use types like:
// - spookytypes.LoggingConfig
// - spookytypes.FactCollection
// - spookytypes.Machine
// etc.
//
// Note: This package uses subpackages to organize types by domain.
// Import specific subpackages for domain-specific types:
//
// import (
//     spookytypesconfig "spooky/internal/types/config"
//     spookytypesfacts "spooky/internal/types/facts"
// )
