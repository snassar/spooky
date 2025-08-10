# Spooky Function Summary

This document provides a comprehensive summary of all functions in the spooky codebase, organized by file for easy code review tracking.

## Internal Package Functions

### `internal/actions/`

#### `internal/actions/actor.go`
- `NewManager(logger)` - Creates new ActionManager with default 30-minute timeout and sequential execution
- `NewManagerWithFacts(logger, factsManager)` - Creates ActionManager with facts manager for enhanced action context
- `Run(ctx, actions, context)` - Orchestrates multiple actions with dependency resolution, returns OrchestrationResult with success/failure status
- `planOrchestration(actions, context)` - Analyzes action dependencies and creates execution plan with parallel/sequential grouping
- `buildDependencyGraph(actions)` - Constructs directed graph of action dependencies for topological sorting
- `optimizeOrchestration(actions, order, graph, context)` - Optimizes execution plan by grouping independent actions for parallel execution
- `runSequential(ctx, runContext)` - Executes actions one after another, returns OrchestrationResult with individual action results
- `runParallel(ctx, runContext)` - Executes independent actions concurrently, returns OrchestrationResult with parallel execution results
- `runAction(ctx, action, context)` - Executes single action (command/script/template), returns RunResult with output and status
- `LoadActions(projectPath)` - Loads actions from project's actions.hcl files, returns ActionCollection with metadata
- `GetAction(name)` - Retrieves action by name from internal registry, returns Action or error if not found
- `ListActions()` - Returns slice of all registered actions with their names and metadata
- `AddAction(name, action)` - Registers new action in internal map, validates action configuration before adding
- `RemoveAction(name)` - Removes action from registry by name, returns error if action doesn't exist
- `ValidateAction(action)` - Validates action configuration (name, type, parameters), returns validation errors
- `SetDefaultTimeout(timeout)` - Sets default timeout for all actions (overrides individual action timeouts)
- `SetDefaultParallel(parallel)` - Sets default parallel execution mode for action collections
- `RegisterCustomValidator(name, validator)` - Registers custom validation function for specific action types
- `Close()` - Cleans up resources, closes connections, and stops any running operations

#### `internal/actions/interfaces.go`
**ActionManager Interface:**
- `LoadActions(projectPath)` - Loads actions from project path, returns ActionCollection or error
- `GetAction(name)` - Gets action by name from registry, returns Action or error if not found
- `ListActions()` - Lists all registered actions, returns slice of Actions
- `AddAction(name, action)` - Adds action to registry with name, returns error if validation fails
- `RemoveAction(name)` - Removes action from registry by name, returns error if not found
- `ExecuteAction(ctx, action, context)` - Executes single action with context, returns ActingSession with results
- `ExecuteActionCollection(ctx, collection, context)` - Executes action collection with context, returns ActingSession
- `PrepareAction(action, context)` - Prepares action for execution by validating and setting up resources
- `PlanAction(action, context)` - Creates execution plan for single action, returns ActionPlan
- `PlanActionCollection(collection, context)` - Creates execution plan for action collection, returns ActionPlan
- `ValidatePlan(plan)` - Validates action execution plan, returns validation errors
- `ValidateAction(action)` - Validates single action configuration, returns validation errors
- `ValidateActionCollection(collection)` - Validates action collection configuration, returns validation errors
- `ValidateActionContext(context)` - Validates action context configuration, returns validation errors
- `MergeActions(actions...)` - Merges multiple actions into collection, returns ActionCollection
- `MergeWithPolicy(existing, new, policy)` - Merges action collections with specified policy, returns merged collection
- `OptimizeAction(action)` - Optimizes action for better performance, returns error if optimization fails
- `OptimizeActionCollection(collection)` - Optimizes action collection for better performance
- `GetPerformanceMetrics(action)` - Gets performance metrics for action, returns PerformanceMetrics
- `SetDefaultTimeout(timeout)` - Sets default timeout for all actions
- `SetDefaultParallel(parallel)` - Sets default parallel execution mode
- `RegisterCustomValidator(name, validator)` - Registers custom validator for specific validation logic
- `Close()` - Closes action manager and cleans up resources

**ActionValidator Interface:**
- `Validate(action)` - Validates single action configuration, returns validation errors
- `ValidateCollection(collection)` - Validates action collection configuration, returns validation errors
- `ValidateContext(context)` - Validates action context configuration, returns validation errors

#### `internal/actions/manager.go`
- `NewManager(logger)` - Creates new ActionManager
- `LoadActions(projectPath)` - Loads actions from project
- `GetAction(name)` - Gets action by name
- `ListActions()` - Lists all actions
- `AddAction(name, action)` - Adds new action
- `RemoveAction(name)` - Removes action
- `ValidateAction(action)` - Validates action
- `SetDefaultTimeout(timeout)` - Sets default timeout
- `SetDefaultParallel(parallel)` - Sets default parallel setting
- `RegisterCustomValidator(name, validator)` - Registers custom validator
- `Close()` - Closes manager

#### `internal/actions/acting/actor.go`
- `NewActor(action, context)` - Creates new actor instance for specific action with initial pending state
- `Execute(ctx, context)` - Executes action based on type (command/script/template), updates result with output and status
- `Prepare()` - Prepares actor for execution by validating action and setting up resources
- `Cancel()` - Cancels actor execution by setting state to cancelled and stopping any running operations
- `GetState()` - Returns current actor state (pending/running/completed/failed/cancelled)
- `GetStatus()` - Returns current actor status with detailed progress information
- `executeCommand(ctx, context, result)` - Executes command action on target machines, captures stdout/stderr
- `executeScript(ctx, context, result)` - Executes script action by uploading and running script file on target machines
- `executeTemplateDeploy(ctx, context, result)` - Deploys template by rendering and copying files to target machines
- `executeTemplateEvaluate(ctx, context, result)` - Evaluates template without deploying, returns rendered content
- `executeTemplateValidate(ctx, context, result)` - Validates template syntax and references without execution
- `executeTemplateCleanup(ctx, context, result)` - Cleans up temporary files and resources after template operations

#### `internal/actions/acting/manager.go`
- `NewManager(executor, sessionManager, resultProcessor, progressTracker)` - Creates acting manager with dependency injection for all components
- `ExecuteAction(ctx, action, context)` - Executes single action with session tracking, returns ActingSession with results
- `ExecuteActionCollection(ctx, collection, context)` - Executes multiple actions with dependency resolution, returns ActingSession with collection results
- `GetSession(sessionID)` - Retrieves acting session by ID, returns ActingSession or error if not found
- `ListSessions()` - Returns slice of all active and completed acting sessions with their metadata
- `executeActionWithRetry(ctx, action, context)` - Executes action with configurable retry logic and exponential backoff

#### `internal/actions/acting/interfaces.go`
**ActingManager Interface:**
- `ExecuteAction(ctx, action, context)` - Executes single action with context, returns ActingSession with results
- `ExecuteActionCollection(ctx, collection, context)` - Executes action collection with context, returns ActingSession
- `PrepareAction(action, context)` - Prepares action for execution by validating and setting up resources
- `CreateActor(action, context)` - Creates actor instance for action execution, returns Actor or error
- `GetActor(action)` - Gets existing actor for action, returns Actor or error if not found
- `GetSession(sessionID)` - Gets acting session by ID, returns ActingSession or error if not found
- `ListSessions()` - Lists all acting sessions, returns slice of ActingSessions
- `CancelSession(sessionID)` - Cancels acting session by ID, returns error if session not found
- `SetDefaultTimeout(timeout)` - Sets default timeout for all acting operations
- `SetDefaultParallel(parallel)` - Sets default parallel execution mode for acting operations
- `SetMaxConcurrent(maxConcurrent)` - Sets maximum concurrent acting operations

**Actor Interface:**
- `Execute(ctx, context)` - Executes action with context, returns ActingResult with execution details
- `Prepare(context)` - Prepares actor for execution, returns error if preparation fails
- `Cancel()` - Cancels actor execution, returns error if cancellation fails
- `GetState()` - Gets current actor state (pending/running/completed/failed/cancelled)
- `GetProgress()` - Gets execution progress as percentage (0.0 to 1.0)
- `GetStatus()` - Gets current actor status with detailed information
- `SetTimeout(timeout)` - Sets timeout for actor execution
- `SetParallel(parallel)` - Sets parallel execution mode for actor

**ActingExecutor Interface:**
- `ExecuteCommand(ctx, command, context)` - Executes command on machines, returns ActingResult
- `ExecuteScript(ctx, script, context)` - Executes script on machines, returns ActingResult
- `ExecuteTemplate(ctx, template, context)` - Executes template on machines, returns ActingResult
- `GetMachine(machineID)` - Gets machine by ID, returns Machine or error if not found
- `ListMachines()` - Lists all available machines, returns slice of Machines
- `ConnectMachine(machineID)` - Connects to machine by ID, returns error if connection fails
- `DisconnectMachine(machineID)` - Disconnects from machine by ID, returns error if disconnection fails
- `SetConnectionTimeout(timeout)` - Sets connection timeout for machine operations
- `SetCommandTimeout(timeout)` - Sets command execution timeout
- `SetRetryAttempts(attempts)` - Sets number of retry attempts for failed operations
- `SetRetryDelay(delay)` - Sets delay between retry attempts

**Machine Interface:**
- `GetID()` - Gets machine ID, returns string
- `GetName()` - Gets machine name, returns string
- `GetHost()` - Gets machine host address, returns string
- `GetUser()` - Gets machine user, returns string
- `GetPort()` - Gets machine port, returns int
- `Connect()` - Connects to machine, returns error if connection fails
- `Disconnect()` - Disconnects from machine, returns error if disconnection fails
- `IsConnected()` - Checks if machine is connected, returns boolean
- `ExecuteCommand(ctx, command)` - Executes command on machine, returns ActingResult
- `ExecuteScript(ctx, script)` - Executes script on machine, returns ActingResult
- `UploadFile(ctx, localPath, remotePath)` - Uploads file to machine, returns error if upload fails
- `DownloadFile(ctx, remotePath, localPath)` - Downloads file from machine, returns error if download fails
- `SetTimeout(timeout)` - Sets timeout for machine operations
- `SetSudo(sudo)` - Sets sudo mode for machine operations
- `SetWorkingDirectory(dir)` - Sets working directory for machine operations
- `SetEnvironment(env)` - Sets environment variables for machine operations

**ActingSessionManager Interface:**
- `CreateSession(actionName)` - Creates new acting session, returns ActingSession or error
- `GetSession(sessionID)` - Gets session by ID, returns ActingSession or error if not found
- `ListSessions()` - Lists all sessions, returns slice of ActingSessions
- `UpdateSession(session)` - Updates session information, returns error if update fails
- `DeleteSession(sessionID)` - Deletes session by ID, returns error if deletion fails
- `StartSession(sessionID)` - Starts session execution, returns error if start fails
- `CompleteSession(sessionID)` - Marks session as completed, returns error if completion fails
- `FailSession(sessionID, err)` - Marks session as failed with error, returns error if operation fails
- `CancelSession(sessionID)` - Cancels session execution, returns error if cancellation fails
- `CleanupExpiredSessions(maxAge)` - Removes expired sessions older than maxAge, returns error if cleanup fails
- `CleanupAllSessions()` - Removes all sessions, returns error if cleanup fails

**ActingResultProcessor Interface:**
- `ProcessResult(result)` - Processes single acting result, returns error if processing fails
- `ProcessResults(results)` - Processes multiple acting results, returns error if processing fails
- `AggregateResults(results)` - Aggregates multiple results into session, returns ActingSession or error
- `CalculateSuccessRate(results)` - Calculates success rate from results, returns float64 percentage
- `ValidateResult(result)` - Validates acting result, returns validation errors
- `ValidateResults(results)` - Validates multiple acting results, returns validation errors
- `TransformResult(result, format)` - Transforms result to specified format, returns transformed data or error
- `TransformResults(results, format)` - Transforms multiple results to specified format, returns transformed data or error

**ActingProgressTracker Interface:**
- `StartTracking(sessionID, totalSteps)` - Starts progress tracking for session, returns error if tracking fails
- `UpdateProgress(sessionID, completedSteps)` - Updates progress for session, returns error if update fails
- `GetProgress(sessionID)` - Gets progress for session, returns float64 percentage or error
- `CompleteTracking(sessionID)` - Completes progress tracking for session, returns error if completion fails
- `GetProgressReport(sessionID)` - Gets detailed progress report for session, returns ProgressReport or error
- `ListProgressReports()` - Lists all progress reports, returns slice of ProgressReports
- `CleanupProgress(sessionID)` - Removes progress tracking for session, returns error if cleanup fails
- `CleanupAllProgress()` - Removes all progress tracking, returns error if cleanup fails

### `internal/cli/`

#### `internal/cli/manager.go`
- `NewManager(config, rootCommand, logger)` - Creates CLI manager with configuration and root command setup
- `ExecuteCommand(args)` - Executes CLI command with provided arguments, returns error if command fails
- `SetRootCommand(command)` - Sets the root command for the CLI application (typically spooky command)
- `GetRootCommand()` - Returns the current root command instance for CLI operations
- `Close()` - Closes CLI manager and cleans up any resources or connections

#### `internal/cli/commands/actions.go`
- `CreateActionsCommand()` - Creates main actions command with all subcommands (list, run, validate)
- `createActionsListCommand()` - Creates subcommand to list all available actions in project
- `createActionsRunCommand()` - Creates subcommand to run actions with dry-run support and parallel execution
- `createActionsValidateCommand()` - Creates subcommand to validate action configurations and dependencies

#### `internal/cli/commands/executor.go`
- `NewActor()` - Creates new command actor for executing CLI commands and subcommands
- `RunCommand(command, args)` - Executes command with arguments, returns error if execution fails
- `RunSubcommand(parent, subcommand, args)` - Executes subcommand under parent command with arguments
- `ValidateCommand(command, args)` - Validates command structure and arguments before execution

#### `internal/cli/commands/manager.go`
- `NewManager(config, builder, executor, logger)` - Creates commands manager with dependency injection for command building and execution
- `RegisterCommand(command)` - Registers new command in internal registry, validates command structure before registration
- `UnregisterCommand(name)` - Removes command from registry by name, returns error if command doesn't exist
- `GetCommand(name)` - Retrieves command by name from registry, returns Command or error if not found
- `ListCommands()` - Returns slice of all registered commands with their names and metadata
- `ValidateCommand(command)` - Validates command structure, flags, and subcommands, returns validation errors
- `Close()` - Closes commands manager and cleans up any registered commands or resources

#### `internal/cli/commands/interfaces.go`
**CommandsManager Interface:**
- `RegisterCommand(command)` - Registers new command in command registry, returns error if validation fails
- `UnregisterCommand(name)` - Removes command from registry by name, returns error if not found
- `GetCommand(name)` - Gets command by name from registry, returns Command or error if not found
- `ListCommands()` - Lists all registered commands, returns slice of Commands
- `InitializeCommands()` - Initializes command system and loads default commands, returns error if initialization fails
- `CreateActionsCommand()` - Creates actions command with subcommands, returns cobra.Command
- `CreateFactsCommand()` - Creates facts command with subcommands, returns cobra.Command
- `CreateMachinesCommand()` - Creates machines command with subcommands, returns cobra.Command
- `CreateProjectCommand()` - Creates project command with subcommands, returns cobra.Command
- `CreateTemplatesCommand()` - Creates templates command with subcommands, returns cobra.Command
- `CreateVariablesCommand()` - Creates variables command with subcommands, returns cobra.Command
- `SetCommandFlags(commandName, flags)` - Sets flags for specific command, returns error if command not found
- `SetCommandExamples(commandName, examples)` - Sets examples for specific command, returns error if command not found
- `ValidateCommand(command)` - Validates command structure and configuration, returns validation errors
- `Close()` - Closes commands manager and cleans up resources

**CommandBuilder Interface:**
- `BuildCommand(command)` - Builds cobra.Command from command definition, returns cobra.Command or error
- `BuildSubcommands(parent, subcommands)` - Builds subcommands for parent command, returns error if build fails
- `ValidateCommandStructure(command)` - Validates command structure and dependencies, returns validation errors

**CommandExecutor Interface:**
- `ExecuteCommand(command, args)` - Executes command with arguments, returns error if execution fails
- `ExecuteSubcommand(parent, subcommand, args)` - Executes subcommand under parent command, returns error if execution fails
- `ValidateExecution(command, args)` - Validates command execution parameters, returns validation errors

#### `internal/cli/factory.go`
- `NewManager(config, logger)` - Creates CLI manager with dependencies
- `NewExecutor()` - Creates command executor

#### `internal/cli/interfaces.go`
**CLIManager Interface:**
- `InitializeCommands()` - Initializes CLI commands and sets up command structure, returns error if initialization fails
- `ExecuteCommand(args)` - Executes CLI command with arguments, returns error if execution fails
- `GetRootCommand()` - Gets root command instance for CLI operations, returns cobra.Command
- `RegisterCommand(command)` - Registers new command in CLI system, returns error if registration fails
- `UnregisterCommand(name)` - Removes command from CLI system by name, returns error if not found
- `GetCommand(name)` - Gets command by name from CLI system, returns Command or error if not found
- `ListCommands()` - Lists all registered commands in CLI system, returns slice of Commands
- `SetGlobalFlags(flags)` - Sets global flags for CLI application, returns error if validation fails
- `SetCommandFlags(commandName, flags)` - Sets flags for specific command, returns error if command not found
- `EnableCompletion(enabled)` - Enables or disables command completion, returns error if operation fails
- `GenerateCompletion(shell)` - Generates completion script for specified shell, returns completion string or error
- `ShowHelp(commandName)` - Shows help for specified command, returns help string or error
- `Close()` - Closes CLI manager and cleans up resources

**CommandsManager Interface:**
- `RegisterCommand(command)` - Registers new command in command registry, returns error if validation fails
- `UnregisterCommand(name)` - Removes command from registry by name, returns error if not found
- `GetCommand(name)` - Gets command by name from registry, returns Command or error if not found
- `ListCommands()` - Lists all registered commands, returns slice of Commands
- `InitializeCommands()` - Initializes command system and loads default commands, returns error if initialization fails

**CompletionManager Interface:**
- `GenerateCompletion(shell)` - Generates completion script for specified shell, returns completion string or error
- `GenerateCompletionFile(shell, outputPath)` - Generates completion script and saves to file, returns error if write fails
- `GetSupportedShells()` - Gets list of supported shells for completion, returns string slice

**HelpManager Interface:**
- `ShowHelp(commandName)` - Shows help for specified command, returns help string or error
- `ShowUsage(commandName)` - Shows usage information for command, returns usage string or error
- `ShowExamples(commandName)` - Shows examples for command, returns examples string or error

**FlagsManager Interface:**
- `SetGlobalFlags(flags)` - Sets global flags for CLI application, returns error if validation fails
- `SetCommandFlags(commandName, flags)` - Sets flags for specific command, returns error if command not found
- `GetGlobalFlags()` - Gets all global flags, returns map of flag names to values
- `GetCommandFlags(commandName)` - Gets flags for specific command, returns map of flag names to values

### `internal/ssh/`

#### `internal/ssh/manager.go`
- `NewManager(config, clientManager, authenticationManager, connectionPoolManager, actingManager, keyManager, logger)` - Creates SSH manager with all required components for connection management
- `Connect(host, config)` - Establishes SSH connection to host using provided configuration, returns SSHConnection or error
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with output and exit code
- `ExecuteScript(connection, script)` - Uploads and executes script on SSH connection, returns CommandResult with script output
- `ExecuteAction(connection, action)` - Executes action (command/script) on SSH connection, returns ActionResult with execution details
- `ExecuteTemplate(connection, template)` - Renders and executes template on SSH connection, returns ActionResult with template results
- `Close()` - Closes SSH manager and all active connections, cleans up resources

#### `internal/ssh/interfaces.go`
**SSHManager Interface:**
- `Connect(host, config)` - Establishes SSH connection to host using config, returns SSHConnection or error
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with output and exit code
- `ExecuteScript(connection, script)` - Uploads and executes script on SSH connection, returns CommandResult with script output
- `CloseConnection(connection)` - Closes SSH connection and cleans up resources
- `GetConnection(host)` - Gets connection from pool for host, returns SSHConnection or error
- `ReturnConnection(connection)` - Returns connection to pool for reuse
- `CloseAllConnections()` - Closes all connections in pool and cleans up resources
- `Authenticate(connection, auth)` - Authenticates SSH connection using auth config, returns error if authentication fails
- `ValidateAuthentication(auth)` - Validates authentication configuration, returns validation errors
- `ExecuteAction(connection, action)` - Executes action on SSH connection, returns ActionResult with execution details
- `ExecuteTemplate(connection, template)` - Renders and executes template on SSH connection, returns ActionResult
- `SetDefaultTimeout(timeout)` - Sets default timeout for SSH operations, returns error if invalid timeout
- `SetMaxConnections(max)` - Sets maximum number of connections in pool, returns error if invalid value
- `EnableConnectionPooling(enabled)` - Enables or disables connection pooling, returns error if operation fails
- `TestConnection(host)` - Tests SSH connection to host, returns error if connection fails
- `GetConnectionStats()` - Gets connection pool statistics, returns ConnectionStats
- `Close()` - Closes SSH manager and all connections

**SSHClient Interface:**
- `Connect(host, config)` - Establishes SSH connection to host using config, returns SSHConnection or error
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with output
- `ExecuteScript(connection, script)` - Uploads and executes script on SSH connection, returns CommandResult
- `CloseConnection(connection)` - Closes SSH connection and cleans up resources

**AuthenticationEngine Interface:**
- `Authenticate(connection, auth)` - Authenticates SSH connection using auth config, returns error if authentication fails
- `ValidateAuthentication(auth)` - Validates authentication configuration, returns validation errors
- `GetSupportedMethods()` - Gets list of supported authentication methods, returns string slice

**ConnectionPool Interface:**
- `GetConnection(host)` - Gets connection from pool for host, returns SSHConnection or error
- `ReturnConnection(connection)` - Returns connection to pool for reuse
- `CloseConnection(connection)` - Closes specific connection and removes from pool
- `CloseAllConnections()` - Closes all connections in pool and cleans up resources
- `GetStats()` - Gets connection pool statistics, returns PoolStats

**ActingEngine Interface:**
- `ExecuteAction(connection, action)` - Executes action on SSH connection, returns ActionResult with execution details
- `ExecuteTemplate(connection, template)` - Renders and executes template on SSH connection, returns ActionResult
- `ExecuteSequential(connection, actions)` - Executes actions sequentially on SSH connection, returns combined ActionResult
- `ExecuteParallel(connection, actions)` - Executes actions in parallel on SSH connection, returns combined ActionResult

**SSHKeyManager Interface:**
- `LoadPrivateKey(path)` - Loads SSH private key from file path, returns SSHKey or error
- `ValidateKeyFile(path)` - Validates SSH key file format and permissions, returns error if invalid

#### `internal/ssh/client/manager.go`
- `NewManager(config, logger)` - Creates SSH client manager with connection pooling and authentication
- `Connect(host, config)` - Establishes SSH connection to host using authentication config, returns SSHConnection
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with stdout/stderr
- `ExecuteScript(connection, script)` - Uploads script to remote host and executes it, returns CommandResult with script output
- `validateExecutionParams(connection, command)` - Validates SSH connection and command parameters before execution
- `Close()` - Closes SSH client manager and terminates all active connections

#### `internal/ssh/client/interfaces.go`
**ClientManager Interface:**
- `Connect(host, config)` - Establishes SSH connection to host using config, returns SSHConnection or error
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with stdout/stderr
- `ExecuteScript(connection, script)` - Uploads script to remote host and executes it, returns CommandResult with script output
- `CloseConnection(connection)` - Closes SSH connection and cleans up resources
- `SetDefaultTimeout(timeout)` - Sets default timeout for SSH operations, returns error if invalid timeout
- `SetMaxRetries(max)` - Sets maximum retry attempts for failed operations, returns error if invalid value
- `EnableHostKeyChecking(enabled)` - Enables or disables host key checking, returns error if operation fails
- `TestConnection(host)` - Tests SSH connection to host, returns error if connection fails
- `GetConnectionInfo(connection)` - Gets detailed connection information, returns ConnectionInfo
- `Close()` - Closes client manager and terminates all active connections

**ConnectionManager Interface:**
- `Connect(host, config)` - Establishes SSH connection to host using config, returns SSHConnection or error
- `CloseConnection(connection)` - Closes SSH connection and cleans up resources
- `TestConnection(connection)` - Tests SSH connection health, returns error if connection is unhealthy
- `GetConnectionState(connection)` - Gets current connection state, returns ConnectionState

**ExecutionManager Interface:**
- `ExecuteCommand(connection, command)` - Executes command on SSH connection, returns CommandResult with output
- `ExecuteScript(connection, script)` - Uploads and executes script on SSH connection, returns CommandResult
- `ExecuteWithTimeout(connection, command, timeout)` - Executes command with specific timeout, returns CommandResult

**HostKeyManager Interface:**
- `GetHostKeyCallback()` - Gets host key callback function for SSH connections, returns HostKeyCallback
- `ValidateHostKey(host, key)` - Validates host key for specific host, returns error if validation fails
- `AddKnownHost(host, key)` - Adds host key to known hosts, returns error if operation fails
- `RemoveKnownHost(host)` - Removes host key from known hosts, returns error if operation fails

#### `internal/ssh/acting/manager.go`
- `NewManager(config, logger)` - Creates acting manager for SSH-based action execution with configuration
- `RunAction(connection, action)` - Runs action on SSH connection, returns ActionResult with execution results
- `RunTemplate(connection, template)` - Renders and runs template on SSH connection, returns ActionResult with template output
- `RunSequential(connection, actions)` - Runs multiple actions sequentially on SSH connection, returns combined ActionResult
- `RunParallel(connection, actions)` - Runs multiple actions in parallel on SSH connection, returns combined ActionResult

#### `internal/ssh/acting/interfaces.go`
**ActingEngine Interface:**
- `RunAction(connection, action)` - Runs action on SSH connection, returns ActionResult with execution results
- `RunTemplate(connection, template)` - Renders and runs template on SSH connection, returns ActionResult with template output
- `RunSequential(connection, actions)` - Runs multiple actions sequentially on SSH connection, returns combined ActionResult
- `RunParallel(connection, actions)` - Runs multiple actions in parallel on SSH connection, returns combined ActionResult

### `internal/project/`

#### `internal/project/manager.go`
- `NewManager(config, logger)` - Creates project manager with configuration and logging setup
- `CreateProject(name, config)` - Creates new project with specified name and configuration, returns Project instance
- `LoadProject(path)` - Loads project from file path, parses project.hcl and related files, returns Project or error
- `ValidateProject(project)` - Validates project configuration against schemas, returns validation errors
- `UpdateProject(existingProject, updates)` - Updates existing project with new configuration, returns updated Project
- `Close()` - Closes project manager and cleans up any loaded projects or resources

#### `internal/project/configuration.go`
- `NewProject()` - Creates new empty project with default configuration values
- `LoadProject(path)` - Loads project from file path, parses HCL configuration, returns Project with resolved settings
- `ValidateProject(project)` - Validates project configuration against schema rules, returns validation errors
- `ResolveProject(project)` - Resolves project configuration by applying defaults and resolving references, returns resolved Project
- `ApplyExecutionSettings(project, resolved)` - Applies execution settings from project to resolved configuration
- `ValidateProjectStructure(project)` - Validates project directory structure against required files and folders
- `Close()` - Closes configuration manager and cleans up any loaded configurations

#### `internal/project/structure.go`
- `NewProjectStructureEngine()` - Creates structure engine
- `ValidateProjectStructure(project)` - Validates project structure
- `CreateProjectStructure(project)` - Creates project structure
- `GenerateProjectFiles(project)` - Generates project files
- `Close()` - Closes structure engine

#### `internal/project/portability.go`
- `NewProjectPortabilityEngine()` - Creates portability engine
- `ValidateProjectPortability(project)` - Validates project portability
- `GenerateEnvironmentAliases(project)` - Generates environment aliases
- `GetPortableName(project, environment)` - Gets portable name
- `CreatePortabilityContext(project, environment)` - Creates portability context
- `MigrateProject(project, targetEnvironment)` - Migrates project
- `GetEnvironmentMappings()` - Gets environment mappings
- `IsPortable(project)` - Checks if project is portable
- `GetPortabilityIssues(project)` - Gets portability issues
- `CreatePortableProject(name, environment, region)` - Creates portable project
- `GetEnvironmentFromPath(projectPath)` - Gets environment from path

#### `internal/project/dependencies.go`
- `NewProjectDependenciesEngine()` - Creates dependencies engine
- `BuildDependencyGraph(projects)` - Builds dependency graph
- `ValidateDependencies(projects)` - Validates dependencies
- `ImportProject(projectPath, importPath, projects)` - Imports project
- `GetDependencyOrder(projects)` - Gets dependency order
- `GetProjectDependencies(projectPath, projects)` - Gets project dependencies
- `GetProjectDependents(projectPath, projects)` - Gets project dependents

### `internal/facts/`

#### `internal/facts/interfaces.go`
- Interface definitions for fact management

#### `internal/facts/format_detector.go`
- `NewFormatDetector()` - Creates format detector for identifying data formats (JSON, HCL, etc.)
- `DetectFormat(r)` - Detects format of data from reader, returns format string (json, hcl, etc.) or error
- `isValidJSON(content)` - Validates JSON content by parsing, returns true if valid JSON structure
- `isValidHCL(content)` - Validates HCL content by parsing, returns true if valid HCL syntax
- `validateAgainstJSONSchema(content)` - Validates JSON content against predefined schema, returns validation errors
- `validateAgainstHCLSchema(content)` - Validates HCL content against predefined schema, returns validation errors

#### `internal/facts/collectors/json/collector.go`
- `NewCollector(sourcePath, logger)` - Creates JSON collector for reading facts from JSON files or directories
- `Collect(server)` - Collects all facts from JSON source for specified server, returns FactCollection
- `CollectSpecific(server, keys)` - Collects specific facts by keys from JSON source, returns FactCollection
- `GetFact(server, key)` - Gets single fact by key from JSON source, returns Fact or error if not found
- `Validate()` - Validates JSON collector configuration and source path accessibility
- `collectFromFile(collection, filePath)` - Collects facts from single JSON file, adds to existing collection
- `collectFromDirectory(collection)` - Collects facts from all JSON files in directory, adds to existing collection
- `convertJSONToFacts(data, source)` - Converts JSON data structure to facts map with source metadata
- `convertValueToFacts(prefix, value, facts, source)` - Recursively converts JSON values to facts with dot notation keys
- `validateJSONFile(filePath)` - Validates JSON file syntax and structure, returns error if invalid

#### `internal/facts/collectors/local/collector.go`
- `NewCollector()` - Creates local collector for gathering system facts from current machine
- `Collect(server)` - Collects all local system facts for specified server, returns FactCollection with comprehensive system data
- `CollectSpecific(server, keys)` - Collects specific facts by keys from local system, returns FactCollection with requested facts
- `GetFact(server, key)` - Gets single fact by key from local system, returns Fact or error if not available
- `Validate()` - Validates local collector can access system information and required files
- `collectSystemFacts(collection)` - Collects basic system facts (hostname, machine ID, FQDN)
- `collectOSFacts(collection)` - Collects operating system facts (name, version, distribution, architecture)
- `collectHardwareFacts(collection)` - Collects hardware facts (CPU cores, model, memory, disk space)
- `collectNetworkFacts(collection)` - Collects network facts (IP addresses, MAC addresses, DNS configuration)
- `collectEnvironmentFacts(collection)` - Collects environment facts (environment variables, user info)
- `collectEnhancedFacts(collection)` - Collects enhanced facts requiring additional system calls or file reads
- `collectSpecificFact(collection, key)` - Collects single fact by key using appropriate collection method
- Various specific fact collection methods (hostname, machineID, OS info, etc.) - Individual methods for specific system information gathering

### `internal/templates/`

#### `internal/templates/manager.go`
- `NewManager(config, engineManager, functionsManager, validationManager, secretsManager, logger)` - Creates template manager with all required components for template processing
- `LoadTemplate(path)` - Loads template from file path, reads content and metadata, returns Template instance
- `RenderTemplate(templateFile, projectPath, additionalData)` - Renders template with project context and additional data, returns rendered string
- `ValidateTemplate(templateFile)` - Validates template syntax and references, returns validation errors
- `ValidateTemplates(projectPath)` - Validates all templates in project directory, returns slice of validation errors
- `NewTemplateContext(projectPath)` - Creates template context with project variables, facts, and functions
- `GetTemplateFunctions(ctx)` - Gets template functions for context, returns FuncMap with available functions
- `RenderTemplateWithTimeout(ctx, templateFile, projectPath, additionalData)` - Renders template with timeout context, returns rendered string or timeout error
- `SetRenderTimeout(timeout)` - Sets maximum time allowed for template rendering operations
- `SetMaxTemplateSize(maxSize)` - Sets maximum allowed template file size in bytes
- `SetDefaultTimeout(timeout)` - Sets default timeout for template operations when not specified
- `RegisterCustomFunction(name, fn)` - Registers custom function for use in templates, validates function signature
- `Close()` - Closes template manager and cleans up all sub-managers and resources

#### `internal/templates/interfaces.go`
**TemplateManager Interface:**
- `LoadTemplate(path)` - Loads template from file path, reads content and metadata, returns Template instance
- `RenderTemplate(templateFile, projectPath, additionalData)` - Renders template with project context and additional data, returns rendered string
- `ValidateTemplate(templateFile)` - Validates template syntax and references, returns validation errors
- `ValidateTemplates(projectPath)` - Validates all templates in project directory, returns slice of validation errors
- `NewTemplateContext(projectPath)` - Creates template context with project variables, facts, and functions
- `GetTemplateFunctions(ctx)` - Gets template functions for context, returns FuncMap with available functions
- `RenderTemplateWithTimeout(ctx, templateFile, projectPath, additionalData)` - Renders template with timeout context, returns rendered string or timeout error
- `SetRenderTimeout(timeout)` - Sets maximum time allowed for template rendering operations
- `SetMaxTemplateSize(maxSize)` - Sets maximum allowed template file size in bytes
- `SetDefaultTimeout(timeout)` - Sets default timeout for template operations when not specified
- `RegisterCustomFunction(name, fn)` - Registers custom function for use in templates, validates function signature
- `Close()` - Closes template manager and cleans up all sub-managers and resources

**TemplateEngine Interface:**
- `ParseTemplate(content, name)` - Parses template content into template.Template, returns parsed template or error
- `RenderTemplate(tmpl, data)` - Renders template with provided data, returns rendered string or error
- `ValidateTemplate(tmpl)` - Validates template syntax and structure, returns validation errors
- `GetTemplateFunctions()` - Gets built-in template functions, returns FuncMap with available functions

**TemplateValidator Interface:**
- `ValidateTemplate(path)` - Validates template file at path, returns validation errors
- `ValidateTemplates(projectPath)` - Validates all templates in project directory, returns slice of validation errors
- `ValidateSyntax(content)` - Validates template syntax in content, returns validation errors
- `ValidateFunctions(tmpl)` - Validates template functions, returns validation errors

**TemplateFunctions Interface:**
- `GetBuiltinFunctions()` - Gets built-in template functions, returns FuncMap with available functions
- `RegisterCustomFunction(name, fn)` - Registers custom function for use in templates, validates function signature
- `ValidateFunction(name, fn)` - Validates function name and signature, returns validation errors

**TemplateSecrets Interface:**
- `EncryptValue(value)` - Encrypts value for use in templates, returns encrypted string or error
- `DecryptValue(encryptedValue)` - Decrypts encrypted value from templates, returns decrypted string or error
- `ValidateEncryptionKey(key)` - Validates encryption key format and accessibility, returns error if invalid

#### `internal/templates/engine/manager.go`
- `NewManager(config, logger)` - Creates template engine manager with configuration and logging
- `ParseTemplate(content, name)` - Parses template content into template.Template, returns parsed template or error
- `RenderTemplate(tmpl, data)` - Renders template with provided data, returns rendered string or error
- `ValidateTemplate(tmpl)` - Validates template syntax and structure, returns validation errors
- `SetDelimiters(left, right)` - Sets custom template delimiters (default: {{ }})
- `SetMaxExecutionTime(timeout)` - Sets maximum time allowed for template execution
- `EnableStrictMode(strict)` - Enables strict mode for template parsing and rendering
- `GetTemplateFunctions()` - Gets built-in template functions, returns FuncMap with available functions
- `Close()` - Closes engine manager and cleans up any resources

#### `internal/templates/engine/interfaces.go`
**EngineManager Interface:**
- `ParseTemplate(content, name)` - Parses template content into template.Template, returns parsed template or error
- `RenderTemplate(tmpl, data)` - Renders template with provided data, returns rendered string or error
- `ValidateTemplate(tmpl)` - Validates template syntax and structure, returns validation errors
- `SetDelimiters(left, right)` - Sets custom template delimiters (default: {{ }})
- `SetMaxExecutionTime(timeout)` - Sets maximum time allowed for template execution
- `EnableStrictMode(strict)` - Enables strict mode for template parsing and rendering
- `GetTemplateFunctions()` - Gets built-in template functions, returns FuncMap with available functions
- `Close()` - Closes engine manager and cleans up any resources

**TemplateParser Interface:**
- `Parse(content, name)` - Parses template content into template.Template, returns parsed template or error
- `ParseFile(path)` - Parses template file at path, returns parsed template or error
- `ValidateSyntax(content)` - Validates template syntax in content, returns validation errors

**TemplateRenderer Interface:**
- `Render(tmpl, data)` - Renders template with provided data, returns rendered string or error
- `RenderWithTimeout(tmpl, data, timeout)` - Renders template with timeout, returns rendered string or timeout error
- `ValidateData(data)` - Validates data structure for template rendering, returns validation errors

#### `internal/templates/functions/manager.go`
- `NewManager(config, logger)` - Creates template functions manager with configuration and logging
- `GetBuiltinFunctions()` - Gets built-in template functions (upper, lower, trim, math ops), returns FuncMap
- `RegisterCustomFunction(name, fn)` - Registers custom function for use in templates, validates function signature
- `ValidateFunction(name, fn)` - Validates function name and signature against template requirements
- `GetFunction(name)` - Gets function by name from built-in or custom functions, returns function or false
- `ListFunctions()` - Lists all available function names (built-in and custom), returns string slice
- `RemoveFunction(name)` - Removes custom function by name, returns error if function doesn't exist
- `EnableBuiltinFunctions(enabled)` - Enables or disables built-in template functions
- `SetFunctionTimeout(timeout)` - Sets timeout for function execution in templates
- `Close()` - Closes functions manager and cleans up any registered functions

#### `internal/templates/functions/interfaces.go`
**FunctionsManager Interface:**
- `GetBuiltinFunctions()` - Gets built-in template functions (upper, lower, trim, math ops), returns FuncMap
- `RegisterCustomFunction(name, fn)` - Registers custom function for use in templates, validates function signature
- `ValidateFunction(name, fn)` - Validates function name and signature against template requirements
- `GetFunction(name)` - Gets function by name from built-in or custom functions, returns function or false
- `ListFunctions()` - Lists all available function names (built-in and custom), returns string slice
- `RemoveFunction(name)` - Removes custom function by name, returns error if function doesn't exist
- `EnableBuiltinFunctions(enabled)` - Enables or disables built-in template functions
- `SetFunctionTimeout(timeout)` - Sets timeout for function execution in templates
- `Close()` - Closes functions manager and cleans up any registered functions

**FunctionValidator Interface:**
- `ValidateFunction(name, fn)` - Validates function name and signature, returns validation errors
- `ValidateFunctionSignature(fn)` - Validates function signature (input/output parameters), returns validation errors
- `ValidateFunctionReturnType(fn)` - Validates function return type compatibility with templates, returns validation errors

#### `internal/templates/functions/validator.go`
- `NewFunctionValidator()` - Creates function validator for template function validation
- `ValidateFunction(name, fn)` - Validates function name and signature, returns validation errors
- `ValidateFunctionSignature(fn)` - Validates function signature (input/output parameters), returns validation errors
- `ValidateFunctionReturnType(fn)` - Validates function return type compatibility with templates, returns validation errors

### `internal/dependency/`

#### `internal/dependency/manager.go`
- `NewDependencyManager(logger)` - Creates dependency manager for tracking variable and action dependencies
- `AddVariable(name, file, line, dependencies)` - Adds variable with its dependencies to dependency graph
- `AddAction(name, file, line, dependencies)` - Adds action with its dependencies to dependency graph
- `ValidateDependencies()` - Validates all dependencies for circular references and missing items, returns validation errors
- `GetResolutionOrder()` - Gets topological sort order for dependency resolution, returns string slice
- `GetDependencyChain(name)` - Gets complete dependency chain for item, returns string slice of dependencies
- `GetDependents(name)` - Gets all items that depend on specified item, returns string slice
- `GetDependencies(name)` - Gets all items that specified item depends on, returns string slice
- `GetDependencyStats()` - Gets dependency statistics (total nodes, edges, circular deps), returns map with stats
- `VisualizeDependencies()` - Creates visual representation of dependency graph, returns string
- `Clear()` - Clears dependency graph and removes all nodes and edges
- `CacheResult(key, value)` - Caches computation result for dependency analysis
- `GetCachedResult(key)` - Gets cached result by key, returns value and existence flag
- `ClearCache()` - Clears all cached dependency analysis results
- `GetNodeInfo(name)` - Gets detailed information about dependency node, returns map with node data
- `ValidateVariableDependencies(name)` - Validates dependencies for specific variable, returns validation errors
- `GetImpactAnalysis(name)` - Gets impact analysis showing what would be affected if item changes, returns analysis map
- `GetDependencyReport()` - Gets comprehensive dependency report with statistics and issues, returns report map

#### `internal/dependency/actions.go`
- `NewActionDependencyManager()` - Creates action dependency manager for tracking action dependencies
- `AddAction(name, dependencies, sourceFile, sourceLine)` - Adds action with its dependencies and source location
- `AddActionCollection(actions)` - Adds multiple actions from collection to dependency graph
- `ValidateDependencies()` - Validates all action dependencies for circular references and missing actions
- `GetExecutionOrder()` - Gets topological sort order for action execution, returns string slice
- `GetParallelGroups()` - Gets groups of actions that can be executed in parallel, returns slice of string slices
- `calculateDependencyLevel(actionName)` - Calculates dependency level (depth) for action in dependency graph
- `GetDependencyChain(actionName)` - Gets complete dependency chain for action, returns string slice
- `GetDependents(actionName)` - Gets all actions that depend on specified action, returns string slice
- `GetDependencies(actionName)` - Gets all actions that specified action depends on, returns string slice
- `GetDependencyStats()` - Gets action dependency statistics (total actions, dependencies, circular deps), returns map
- `VisualizeDependencies()` - Creates visual representation of action dependency graph, returns string
- `detectMissingReferences()` - Detects actions referenced in dependencies but not defined, returns string slice
- `detectSelfReferences()` - Detects actions that depend on themselves, returns string slice
- `detectOrphanedActions()` - Detects actions with no dependencies and no dependents, returns string slice
- `allDependenciesProcessed(actionName, processed)` - Checks if all dependencies for action have been processed

### `internal/schemas/`

#### `internal/schemas/project.go`
- `NewProjectValidator()` - Creates project validator for validating project configurations
- `ValidateProject(project)` - Validates complete project configuration against schema rules, returns validation errors
- `ValidateProjectExecution(execution)` - Validates project execution settings (timeout, parallel, etc.), returns validation errors
- `ValidateProjectStructure(project)` - Validates project directory structure and required files, returns validation errors
- `ValidateProjectConfiguration(project)` - Validates project configuration format and required fields, returns validation errors

#### `internal/schemas/evolution.go`
- `NewSchemaEvolutionManager()` - Creates schema evolution manager for managing schema migrations
- `ExecuteSchemaEvolutionWorkflow(workflow)` - Executes schema evolution workflow with multiple steps, returns evolution results
- `ValidateSchemaEvolution(workflow)` - Validates schema evolution workflow before execution, returns validation errors
- `GetEvolutionSteps(workflow)` - Gets list of evolution steps from workflow, returns step definitions
- `ApplyEvolutionStep(step)` - Applies single evolution step to schema, returns step result

### `internal/secrets/`

#### `internal/secrets/age/client.go`
- `NewClient()` - Creates age encryption client for handling age-based encryption/decryption
- `Encrypt(data, recipients)` - Encrypts data with age encryption using recipient public keys, returns EncryptedValue
- `Decrypt(encrypted, identityStr)` - Decrypts age-encrypted data using identity private key, returns decrypted bytes
- `GenerateKey()` - Generates new age key pair (identity and recipient), returns identity and recipient strings
- `ParseRecipient(recipient)` - Parses and validates age recipient public key format, returns error if invalid
- `ParseIdentity(identity)` - Parses and validates age identity private key format, returns error if invalid
- `EncryptStream(input, output, recipients)` - Encrypts data stream with age encryption, writes to output stream
- `DecryptStream(input, output, identityStr)` - Decrypts age-encrypted stream, writes decrypted data to output stream

#### `internal/secrets/manager.go`
- `NewManager(config, logger)` - Creates secrets manager with age encryption configuration
- `Encrypt(data, recipients)` - Encrypts data using age encryption with specified recipients, returns EncryptedValue
- `Decrypt(encrypted, identityFile)` - Decrypts age-encrypted data using identity from file, returns decrypted bytes
- `GenerateKey()` - Generates new age key pair and saves to configured location, returns key information
- `ValidateKey(key)` - Validates age key format and accessibility, returns validation errors
- `GetStatus()` - Gets secrets manager status (enabled/disabled, key availability), returns status map
- `Close()` - Closes secrets manager and cleans up any encryption resources

### `internal/storage/`

#### `internal/storage/storage.go`
- `matchesQuery(collection, query)` - Matches fact collection against query criteria, returns true if collection matches
- `NewFactStorage(opts)` - Creates new fact storage with specified options (BadgerDB, JSON, HCL), returns FactStorage
- `OpenFactStorage(opts)` - Opens existing fact storage with specified options, returns FactStorage or error
- `NewFactStorageReadOnly(opts)` - Creates read-only fact storage for querying without modification, returns FactStorage

### `internal/logging/`

#### `internal/logging/fields.go`
- `String(key, value)` - Creates structured logging field for string values
- `Int(key, value)` - Creates structured logging field for integer values
- `Int64(key, value)` - Creates structured logging field for 64-bit integer values
- `Float64(key, value)` - Creates structured logging field for 64-bit float values
- `Bool(key, value)` - Creates structured logging field for boolean values
- `Error(err)` - Creates structured logging field for error values with error message
- `Duration(key, durationMs)` - Creates structured logging field for duration values in milliseconds
- `RequestID(id)` - Creates structured logging field for request ID tracking
- `Server(name)` - Creates structured logging field for server name identification
- `Action(name)` - Creates structured logging field for action name identification
- `Host(host)` - Creates structured logging field for host identification
- `Port(port)` - Creates structured logging field for port number
- `StringSlice(key, value)` - Creates structured logging field for string slice values

### `internal/coordinator/`

#### `internal/coordinator/manager.go`
- `NewManager(config, logger)` - Creates coordinator manager for orchestrating system operations
- `ExecuteAction(action, context)` - Executes action through appropriate subsystem (SSH, local, etc.), returns execution result
- `GetStatus()` - Gets coordinator status including active operations and system health, returns status map
- `Close()` - Closes coordinator and stops all active operations, cleans up resources

#### `internal/coordinator/actions.go`
- `NewActionsManager()` - Creates actions manager for coordinating action execution
- `ExecuteAction(action, context)` - Executes single action with context, returns execution result
- `ExecuteActionCollection(collection, context)` - Executes collection of actions with dependency resolution, returns collection result
- `GetAction(name)` - Gets action by name from coordinator's action registry, returns Action or error
- `ListActions()` - Lists all actions available in coordinator's registry, returns slice of Action names
- `ValidateAction(action)` - Validates action configuration and dependencies, returns validation errors

### `internal/machines/`

#### `internal/machines/connectivity/manager.go`
- `NewManager(config, logger)` - Creates connectivity manager for testing machine connectivity
- `TestConnectivity(machine)` - Tests connectivity to machine using SSH, returns connectivity result with details
- `PingMachine(machine)` - Pings machine to check basic network connectivity, returns ping result
- `GetConnectionStatus(machine)` - Gets detailed connection status for machine, returns status with connection details
- `Close()` - Closes connectivity manager and terminates any active connections

### `internal/variables/`

#### `internal/variables/manager.go`
- `NewManager(config, logger)` - Creates variables manager for handling project variables
- `LoadVariables(projectPath)` - Loads variables from project's variables.hcl and variables/ directory, returns VariableCollection
- `GetVariable(name)` - Gets variable by name from loaded variables, returns Variable or error if not found
- `SetVariable(name, value)` - Sets variable value in memory, validates value before setting
- `ListVariables()` - Lists all loaded variables with their names and metadata, returns slice of Variable names
- `ValidateVariables()` - Validates all loaded variables against schemas and dependencies, returns validation errors
- `Close()` - Closes variables manager and cleans up any loaded variables or resources

### `internal/interfaces/`

#### `internal/interfaces/actions.go`
- Interface definitions for action operations

#### `internal/interfaces/manager.go`
- Interface definitions for manager operations

#### `internal/interfaces/contexts.go`
- Interface definitions for context operations

### `internal/types/`

#### `internal/types/types.go`
- Type definitions and aliases for the entire codebase

## Summary

This document provides a comprehensive overview of all functions in the spooky codebase, organized by package and file. Each function is listed with a brief description of its purpose, making it easy to track during code reviews and understand the overall architecture of the system.

The functions are grouped by their primary responsibility:
- **Actions**: Action management and orchestration
- **CLI**: Command-line interface operations
- **SSH**: SSH connection and execution
- **Project**: Project configuration and management
- **Facts**: Fact collection and management
- **Templates**: Template rendering and management
- **Dependency**: Dependency resolution and management
- **Schemas**: Schema validation and evolution
- **Secrets**: Secret encryption and management
- **Storage**: Data storage operations
- **Logging**: Logging and field management
- **Coordinator**: System coordination
- **Machines**: Machine connectivity
- **Variables**: Variable management
- **Interfaces**: Interface definitions
- **Types**: Type definitions
