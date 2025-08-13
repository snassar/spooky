// Package types provides unified access to all type definitions in the spooky codebase.
// This package serves as the central repository for all type definitions and re-exports
// types from domain-specific subpackages for consistent access patterns.
package types

import (
	spookytypesactions "spooky/internal/types/actions"
	spookytypescli "spooky/internal/types/cli"
	spookytypescommon "spooky/internal/types/common"
	spookytypesconfig "spooky/internal/types/config"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesmachines "spooky/internal/types/machines"
	spookytypesproject "spooky/internal/types/project"
	spookytypesschemas "spooky/internal/types/schemas"
	spookytypesssh "spooky/internal/types/ssh"
	spookytypestemplates "spooky/internal/types/templates"
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

// SSH types
type Connection = spookytypesssh.Connection
type ConnectionStatus = spookytypesssh.ConnectionStatus
type AuthMethod = spookytypesssh.AuthMethod
type ConnectionPool = spookytypesssh.ConnectionPool
type ConnectionFactory = spookytypesssh.ConnectionFactory
type ConnectionRequest = spookytypesssh.ConnectionRequest
type ConnectionResult = spookytypesssh.ConnectionResult

type Client = spookytypesssh.Client
type ClientConfig = spookytypesssh.ClientConfig
type ClientStatus = spookytypesssh.ClientStatus
type Session = spookytypesssh.Session
type SessionStatus = spookytypesssh.SessionStatus
type PtyConfig = spookytypesssh.PtyConfig
type SSHCommand = spookytypesssh.Command
type SSHCommandResult = spookytypesssh.CommandResult
type FileTransfer = spookytypesssh.FileTransfer
type TransferDirection = spookytypesssh.TransferDirection
type TransferMode = spookytypesssh.TransferMode
type FileTransferResult = spookytypesssh.FileTransferResult

type Authentication = spookytypesssh.Authentication
type KeyType = spookytypesssh.KeyType
type Key = spookytypesssh.Key
type KeyPair = spookytypesssh.KeyPair
type HostKey = spookytypesssh.HostKey
type TrustLevel = spookytypesssh.TrustLevel
type KnownHosts = spookytypesssh.KnownHosts
type AuthenticationResult = spookytypesssh.AuthenticationResult

type ActingSession = spookytypesssh.ActingSession
type ActingSessionStatus = spookytypesssh.ActingSessionStatus
type ActingMode = spookytypesssh.ActingMode
type ActingCommand = spookytypesssh.ActingCommand
type ActingCommandResult = spookytypesssh.ActingCommandResult
type ActingBatch = spookytypesssh.ActingBatch
type ActingBatchStatus = spookytypesssh.ActingBatchStatus
type ActingBatchResult = spookytypesssh.ActingBatchResult
type ActingScript = spookytypesssh.ActingScript
type ActingScriptResult = spookytypesssh.ActingScriptResult

type SSHError = spookytypesssh.SSHError
type SSHErrorType = spookytypesssh.SSHErrorType
type ConnectionError = spookytypesssh.ConnectionError
type AuthenticationError = spookytypesssh.AuthenticationError
type HostKeyError = spookytypesssh.HostKeyError
type SSHCommandError = spookytypesssh.CommandError
type SessionError = spookytypesssh.SessionError
type FileTransferError = spookytypesssh.FileTransferError
type ValidationError = spookytypesssh.ValidationError
type ValidationType = spookytypesssh.ValidationType
type ValidationLevel = spookytypesssh.ValidationLevel
type ConfigurationError = spookytypesssh.ConfigurationError

// SSH constants
const (
	// Connection status constants
	ConnectionStatusDisconnected = spookytypesssh.ConnectionStatusDisconnected
	ConnectionStatusConnecting   = spookytypesssh.ConnectionStatusConnecting
	ConnectionStatusConnected    = spookytypesssh.ConnectionStatusConnected
	ConnectionStatusFailed       = spookytypesssh.ConnectionStatusFailed
	ConnectionStatusTimeout      = spookytypesssh.ConnectionStatusTimeout
	ConnectionStatusClosed       = spookytypesssh.ConnectionStatusClosed

	// Auth method constants
	AuthMethodPassword  = spookytypesssh.AuthMethodPassword
	AuthMethodPublicKey = spookytypesssh.AuthMethodPublicKey
	AuthMethodKeyboard  = spookytypesssh.AuthMethodKeyboard
	AuthMethodAgent     = spookytypesssh.AuthMethodAgent
	AuthMethodNone      = spookytypesssh.AuthMethodNone

	// Client status constants
	ClientStatusInitialized  = spookytypesssh.ClientStatusInitialized
	ClientStatusConnecting   = spookytypesssh.ClientStatusConnecting
	ClientStatusConnected    = spookytypesssh.ClientStatusConnected
	ClientStatusDisconnected = spookytypesssh.ClientStatusDisconnected
	ClientStatusError        = spookytypesssh.ClientStatusError
	ClientStatusClosed       = spookytypesssh.ClientStatusClosed

	// Session status constants
	SessionStatusCreated   = spookytypesssh.SessionStatusCreated
	SessionStatusStarting  = spookytypesssh.SessionStatusStarting
	SessionStatusActive    = spookytypesssh.SessionStatusActive
	SessionStatusRunning   = spookytypesssh.SessionStatusRunning
	SessionStatusCompleted = spookytypesssh.SessionStatusCompleted
	SessionStatusFailed    = spookytypesssh.SessionStatusFailed
	SessionStatusClosed    = spookytypesssh.SessionStatusClosed

	// Acting session status constants
	ActingSessionStatusCreated   = spookytypesssh.ActingSessionStatusCreated
	ActingSessionStatusStarting  = spookytypesssh.ActingSessionStatusStarting
	ActingSessionStatusActive    = spookytypesssh.ActingSessionStatusActive
	ActingSessionStatusRunning   = spookytypesssh.ActingSessionStatusRunning
	ActingSessionStatusCompleted = spookytypesssh.ActingSessionStatusCompleted
	ActingSessionStatusFailed    = spookytypesssh.ActingSessionStatusFailed
	ActingSessionStatusClosed    = spookytypesssh.ActingSessionStatusClosed

	// Acting mode constants
	ActingModeInteractive = spookytypesssh.ActingModeInteractive
	ActingModeBatch       = spookytypesssh.ActingModeBatch
	ActingModeScript      = spookytypesssh.ActingModeScript
	ActingModeCommand     = spookytypesssh.ActingModeCommand

	// Acting batch status constants
	ActingBatchStatusPending   = spookytypesssh.ActingBatchStatusPending
	ActingBatchStatusStarting  = spookytypesssh.ActingBatchStatusStarting
	ActingBatchStatusRunning   = spookytypesssh.ActingBatchStatusRunning
	ActingBatchStatusCompleted = spookytypesssh.ActingBatchStatusCompleted
	ActingBatchStatusFailed    = spookytypesssh.ActingBatchStatusFailed
	ActingBatchStatusCancelled = spookytypesssh.ActingBatchStatusCancelled

	// Transfer direction constants
	TransferDirectionUpload   = spookytypesssh.TransferDirectionUpload
	TransferDirectionDownload = spookytypesssh.TransferDirectionDownload

	// Transfer mode constants
	TransferModeSCP  = spookytypesssh.TransferModeSCP
	TransferModeSFTP = spookytypesssh.TransferModeSFTP

	// Key type constants
	KeyTypeRSA     = spookytypesssh.KeyTypeRSA
	KeyTypeDSA     = spookytypesssh.KeyTypeDSA
	KeyTypeECDSA   = spookytypesssh.KeyTypeECDSA
	KeyTypeED25519 = spookytypesssh.KeyTypeED25519

	// Trust level constants
	TrustLevelUnknown   = spookytypesssh.TrustLevelUnknown
	TrustLevelUntrusted = spookytypesssh.TrustLevelUntrusted
	TrustLevelTrusted   = spookytypesssh.TrustLevelTrusted
	TrustLevelVerified  = spookytypesssh.TrustLevelVerified

	// SSH error type constants
	SSHErrorTypeConnection     = spookytypesssh.SSHErrorTypeConnection
	SSHErrorTypeAuthentication = spookytypesssh.SSHErrorTypeAuthentication
	SSHErrorTypeAuthorization  = spookytypesssh.SSHErrorTypeAuthorization
	SSHErrorTypeTimeout        = spookytypesssh.SSHErrorTypeTimeout
	SSHErrorTypeProtocol       = spookytypesssh.SSHErrorTypeProtocol
	SSHErrorTypeHostKey        = spookytypesssh.SSHErrorTypeHostKey
	SSHErrorTypeCommand        = spookytypesssh.SSHErrorTypeCommand
	SSHErrorTypeSession        = spookytypesssh.SSHErrorTypeSession
	SSHErrorTypeFileTransfer   = spookytypesssh.SSHErrorTypeFileTransfer
	SSHErrorTypeValidation     = spookytypesssh.SSHErrorTypeValidation
	SSHErrorTypeConfiguration  = spookytypesssh.SSHErrorTypeConfiguration
	SSHErrorTypeUnknown        = spookytypesssh.SSHErrorTypeUnknown

	// Validation type constants
	ValidationTypeRequired    = spookytypesssh.ValidationTypeRequired
	ValidationTypeFormat      = spookytypesssh.ValidationTypeFormat
	ValidationTypeRange       = spookytypesssh.ValidationTypeRange
	ValidationTypeLength      = spookytypesssh.ValidationTypeLength
	ValidationTypePattern     = spookytypesssh.ValidationTypePattern
	ValidationTypeEnum        = spookytypesssh.ValidationTypeEnum
	ValidationTypeCustom      = spookytypesssh.ValidationTypeCustom
	ValidationTypeDependency  = spookytypesssh.ValidationTypeDependency
	ValidationTypeConsistency = spookytypesssh.ValidationTypeConsistency

	// Validation level constants
	ValidationLevelError   = spookytypesssh.ValidationLevelError
	ValidationLevelWarning = spookytypesssh.ValidationLevelWarning
	ValidationLevelInfo    = spookytypesssh.ValidationLevelInfo
)

// Facts types
type Facts = spookytypesfacts.Facts
type SystemFacts = spookytypesfacts.SystemFacts
type OSFacts = spookytypesfacts.OSFacts
type HardwareFacts = spookytypesfacts.HardwareFacts
type CPUFacts = spookytypesfacts.CPUFacts
type CPUTimes = spookytypesfacts.CPUTimes
type CPUCoreDetail = spookytypesfacts.CPUCoreDetail
type MemoryFacts = spookytypesfacts.MemoryFacts
type SwapFacts = spookytypesfacts.SwapFacts
type VirtualMemoryFacts = spookytypesfacts.VirtualMemoryFacts
type DiskFacts = spookytypesfacts.DiskFacts
type DiskIOCounters = spookytypesfacts.DiskIOCounters
type PartitionFacts = spookytypesfacts.PartitionFacts
type DiskIOFacts = spookytypesfacts.DiskIOFacts
type NetworkFacts = spookytypesfacts.NetworkFacts
type NetworkInterface = spookytypesfacts.NetworkInterface
type NetworkProtocols = spookytypesfacts.NetworkProtocols
type TCPStats = spookytypesfacts.TCPStats
type UDPStats = spookytypesfacts.UDPStats
type ICMPStats = spookytypesfacts.ICMPStats
type NetworkConnection = spookytypesfacts.NetworkConnection
type NetworkInterfaceStats = spookytypesfacts.NetworkInterfaceStats
type NetfilterConntrack = spookytypesfacts.NetfilterConntrack
type LoadAverageFacts = spookytypesfacts.LoadAverageFacts
type ProcessFacts = spookytypesfacts.ProcessFacts
type ProcessInfo = spookytypesfacts.ProcessInfo
type DetailedProcessInfo = spookytypesfacts.DetailedProcessInfo
type ProcessMemoryInfo = spookytypesfacts.ProcessMemoryInfo
type ProcessMemoryMap = spookytypesfacts.ProcessMemoryMap
type ProcessOpenFile = spookytypesfacts.ProcessOpenFile
type ProcessIOCounters = spookytypesfacts.ProcessIOCounters
type ProcessPageFaults = spookytypesfacts.ProcessPageFaults
type EnhancedFacts = spookytypesfacts.EnhancedFacts
type VirtualizationFacts = spookytypesfacts.VirtualizationFacts
type PackageManagerFacts = spookytypesfacts.PackageManagerFacts
type ServiceManagerFacts = spookytypesfacts.ServiceManagerFacts
type SELinuxFacts = spookytypesfacts.SELinuxFacts
type SSHKeysFacts = spookytypesfacts.SSHKeysFacts
type BIOSFacts = spookytypesfacts.BIOSFacts
type SensorsFacts = spookytypesfacts.SensorsFacts
type TemperatureSensors = spookytypesfacts.TemperatureSensors
type FanSensors = spookytypesfacts.FanSensors
type DockerFacts = spookytypesfacts.DockerFacts
type DockerCgroupsCPU = spookytypesfacts.DockerCgroupsCPU
type DockerCgroupsMemory = spookytypesfacts.DockerCgroupsMemory
type ApplicationFacts = spookytypesfacts.ApplicationFacts
type ApplicationVersions = spookytypesfacts.ApplicationVersions
type ApplicationConfig = spookytypesfacts.ApplicationConfig
type ConfigPaths = spookytypesfacts.ConfigPaths
type LogPaths = spookytypesfacts.LogPaths
type DeploymentFacts = spookytypesfacts.DeploymentFacts
type DeploymentState = spookytypesfacts.DeploymentState
type DeploymentInfo = spookytypesfacts.DeploymentInfo
type DeploymentServices = spookytypesfacts.DeploymentServices
type EnvironmentFacts = spookytypesfacts.EnvironmentFacts
type EnvironmentVariables = spookytypesfacts.EnvironmentVariables
type InfrastructureFacts = spookytypesfacts.InfrastructureFacts
type MonitoringFacts = spookytypesfacts.MonitoringFacts
type MonitoringEndpoints = spookytypesfacts.MonitoringEndpoints
type HealthCheckFacts = spookytypesfacts.HealthCheckFacts

// Actions types
type Action = spookytypesactions.Action
type ActionCollection = spookytypesactions.ActionCollection
type ActionPlan = spookytypesactions.ActionPlan
type ActionDependency = spookytypesactions.ActionDependency
type ActionExecution = spookytypesactions.ActionExecution
type ActionValidation = spookytypesactions.ActionValidation
type ActionMetrics = spookytypesactions.ActionMetrics
type ActionExecutionContext = spookytypesactions.ActionExecutionContext
type ActingResult = spookytypesactions.ActingResult
type TemplateConfig = spookytypesactions.TemplateConfig
type FileCopyConfig = spookytypesactions.FileCopyConfig
type ServiceConfig = spookytypesactions.ServiceConfig
type ResourceLimits = spookytypesactions.ResourceLimits

// Template types
type Template = spookytypestemplates.Template
type TemplateMetadata = spookytypestemplates.TemplateMetadata
type TemplateContext = spookytypestemplates.TemplateContext
type TemplateEngine = spookytypestemplates.TemplateEngine
type TemplateManager = spookytypestemplates.TemplateManager
type TemplateError = spookytypestemplates.TemplateError
type TemplateValidationResult = spookytypestemplates.TemplateValidationResult

// Facts types
type FactCollection = spookytypesfacts.FactCollection
type FactStorage = spookytypesfacts.FactStorage
type FactCollector = spookytypesfacts.FactCollector
type FactManager = spookytypesfacts.FactManager
type FactCollectionOptions = spookytypesfacts.FactCollectionOptions
type FactExportOptions = spookytypesfacts.FactExportOptions

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
