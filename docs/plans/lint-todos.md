# Linting TODOs - AI-Optimized Report

This file contains linting issues captured by running `golangci-lint` on each Go file individually.
**Optimized for AI consumption** with categorization, fix suggestions, and progress tracking.

**Generated:** 2025-08-15 08:31:58
**Config:** .golangci.yml
**Fast Mode:** true

## 📊 Summary Statistics

- **Total files processed:** 89
- **Files with issues:** 39
- **Total issues:** 139

### Priority Breakdown

- **🔴 High Priority (Build-breaking):** 69 issues (49.6%)
- **🟡 Medium Priority (Code quality):** 31 issues (22.3%)
- **🟢 Low Priority (Style):** 39 issues (28.1%)

## 🏷️ Issue Categories

### typecheck (69 issues)

**Description:** Build-breaking errors that prevent compilation

**Common Fix Patterns:**
- Add missing import: `import "package/path"`
- Define missing variable: `var variableName Type`
- Fix function call: use correct function name
- Add missing package initialization

**Examples:**
- `cmd/facts.go:344:2`: undefined: RootCmd
- `internal/cli/commands/integrations_test.go:29:9`: undefined: NewIntegrationsCommand
- `internal/cli/commands/integrations_test.go:57:9`: undefined: NewIntegrationsCommand

### revive (30 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/cli/commands/integrations.go:25:36`: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- `internal/cli/commands/integrations.go:57:61`: unused-parameter: parameter 'args' seems to be unused, consider removing or renaming it as _
- `internal/machines/integration.go:416:74`: unused-parameter: parameter 'filter' seems to be unused, consider removing or renaming it as _

### gocritic (29 issues)

**Description:** Code quality and style suggestions

**Common Fix Patterns:**
- Use pointer for large structs: `func process(config *Config)`
- Simplify boolean expressions
- Use consistent naming conventions
- Optimize string concatenation

**Examples:**
- `internal/actions/manager_test.go:390:2`: rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- `internal/actions/manager_test.go:402:2`: rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- `internal/actions/validator.go:76:55`: hugeParam: action is heavy (584 bytes); consider passing it by pointer

### gocyclo (6 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/facts/hcl_parser.go:143:1`: cyclomatic complexity 24 of func `(*HCLParser).parseHostBlock` is high (> 20)
- `internal/machines/manager.go:116:1`: cyclomatic complexity 34 of func `(*Manager).ExportMachines` is high (> 20)
- `internal/machines/manager.go:804:1`: cyclomatic complexity 27 of func `(*Manager).applyMachineFilter` is high (> 20)

### unused (2 issues)

**Description:** Unused variables, imports, or functions

**Common Fix Patterns:**
- Remove unused imports: `import _ "package"`
- Remove unused variables
- Remove unused functions
- Use variables or remove them

**Examples:**
- `cmd/completion.go:11:5`: var `completionCmd` is unused
- `internal/machines/manager.go:549:6`: func `getMachineHostnames` is unused

### gosec (2 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/machines/integration.go:192:12`: G306: Expect WriteFile permissions to be 0600 or less
- `internal/machines/manager.go:255:12`: G306: Expect WriteFile permissions to be 0600 or less

### gosimple (1 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/machines/manager.go:770:5`: S1009: should omit nil check; len() for nil maps is defined as zero

## 📁 File-by-File Breakdown

### cmd/completion.go

**Issues found:** 1

#### 🟢 Low Priority Issues

- **unused** (Low): var `completionCmd` is unused

```bash
# TODO: Fix linting issues in cmd/completion.go
# Run: golangci-lint run --config=.golangci.yml cmd/completion.go
```

### cmd/facts.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: RootCmd

```bash
# TODO: Fix linting issues in cmd/facts.go
# Run: golangci-lint run --config=.golangci.yml cmd/facts.go
```

### internal/actions/manager_test.go

**Issues found:** 2

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)

```bash
# TODO: Fix linting issues in internal/actions/manager_test.go
# Run: golangci-lint run --config=.golangci.yml internal/actions/manager_test.go
```

### internal/actions/validator.go

**Issues found:** 1

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): hugeParam: action is heavy (584 bytes); consider passing it by pointer

```bash
# TODO: Fix linting issues in internal/actions/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/actions/validator.go
```

### internal/cli/commands/integrations.go

**Issues found:** 3

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): rangeValCopy: each iteration copies 248 bytes (consider pointers or indexing)

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'args' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/cli/commands/integrations.go
# Run: golangci-lint run --config=.golangci.yml internal/cli/commands/integrations.go
```

### internal/cli/commands/integrations_test.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewIntegrationsCommand
- **typecheck** (High): undefined: NewIntegrationsCommand

```bash
# TODO: Fix linting issues in internal/cli/commands/integrations_test.go
# Run: golangci-lint run --config=.golangci.yml internal/cli/commands/integrations_test.go
```

### internal/config/auto_setup.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: LoggingConfigManager
- **typecheck** (High): undefined: NewLoggingConfigManager

```bash
# TODO: Fix linting issues in internal/config/auto_setup.go
# Run: golangci-lint run --config=.golangci.yml internal/config/auto_setup.go
```

### internal/facts/collector.go

**Issues found:** 4

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: NewHCLParser
- **typecheck** (High): undefined: NewHCLParser

```bash
# TODO: Fix linting issues in internal/facts/collector.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/collector.go
```

### internal/facts/hcl_parser.go

**Issues found:** 1

#### 🟢 Low Priority Issues

- **gocyclo** (Low): cyclomatic complexity 24 of func `(*HCLParser).parseHostBlock` is high (> 20)

```bash
# TODO: Fix linting issues in internal/facts/hcl_parser.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/hcl_parser.go
```

### internal/facts/integration.go

**Issues found:** 4

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Manager
- **typecheck** (High): undefined: Manager
- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: FactCollection

```bash
# TODO: Fix linting issues in internal/facts/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/integration.go
```

### internal/facts/integration_test.go

**Issues found:** 17

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: MockFactCollector
- **typecheck** (High): undefined: MockSchemaValidator
- **typecheck** (High): undefined: MockLogger
- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewIntegration
- **typecheck** (High): undefined: MockFactCollector
- **typecheck** (High): undefined: MockSchemaValidator
- **typecheck** (High): undefined: MockLogger
- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewIntegration
- **typecheck** (High): undefined: MockFactCollector
- **typecheck** (High): undefined: MockSchemaValidator
- **typecheck** (High): undefined: MockLogger
- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewIntegration
- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: FactCollection

```bash
# TODO: Fix linting issues in internal/facts/integration_test.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/integration_test.go
```

### internal/facts/manager_test.go

**Issues found:** 8

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: FactCollection
- **typecheck** (High): undefined: Manager
- **typecheck** (High): undefined: NewMockFactCollector
- **typecheck** (High): undefined: MockSchemaValidator
- **typecheck** (High): undefined: MockLogger
- **typecheck** (High): undefined: NewManager

```bash
# TODO: Fix linting issues in internal/facts/manager_test.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/manager_test.go
```

### internal/machines/integration.go

**Issues found:** 9

#### 🟡 Medium Priority Issues

- **gosec** (Medium): G306: Expect WriteFile permissions to be 0600 or less
- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): sprintfQuotedString: use %q instead of "%s" for quoted strings
- **gocritic** (Medium): sprintfQuotedString: use %q instead of "%s" for quoted strings
- **gocritic** (Medium): sprintfQuotedString: use %q instead of "%s" for quoted strings

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'filter' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/machines/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/integration.go
```

### internal/machines/loader.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: getMachineHostnames

```bash
# TODO: Fix linting issues in internal/machines/loader.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/loader.go
```

### internal/machines/manager.go

**Issues found:** 13

#### 🟡 Medium Priority Issues

- **gosec** (Medium): G306: Expect WriteFile permissions to be 0600 or less
- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): ifElseChain: rewrite if-else to switch statement
- **gocritic** (Medium): nestingReduce: invert if cond, replace body with `continue`, move old body after the statement
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)

#### 🟢 Low Priority Issues

- **gocyclo** (Low): cyclomatic complexity 34 of func `(*Manager).ExportMachines` is high (> 20)
- **gocyclo** (Low): cyclomatic complexity 27 of func `(*Manager).applyMachineFilter` is high (> 20)
- **gocyclo** (Low): cyclomatic complexity 22 of func `(*Manager).validateMachineCollection` is high (> 20)
- **unused** (Low): func `getMachineHostnames` is unused
- **gosimple** (Low): S1009: should omit nil check; len() for nil maps is defined as zero

```bash
# TODO: Fix linting issues in internal/machines/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/manager.go
```

### internal/machines/validator.go

**Issues found:** 5

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): regexpSimplify: can re-write `^[0-9]+(Gbps|Mbps)$` as `^\d+(Gbps|Mbps)$`

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/machines/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/validator.go
```

### internal/project/loader.go

**Issues found:** 3

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'projectPath' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/project/loader.go
# Run: golangci-lint run --config=.golangci.yml internal/project/loader.go
```

### internal/project/manager.go

**Issues found:** 6

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/project/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/project/manager.go
```

### internal/project/validator.go

**Issues found:** 2

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/project/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/project/validator.go
```

### internal/schemas/evolution_manager.go

**Issues found:** 3

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): rangeValCopy: each iteration copies 136 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 136 bytes (consider pointers or indexing)

#### 🟢 Low Priority Issues

- **revive** (Low): redefines-builtin-id: redefinition of the built-in type error

```bash
# TODO: Fix linting issues in internal/schemas/evolution_manager.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/evolution_manager.go
```

### internal/schemas/manager_test.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewManager

```bash
# TODO: Fix linting issues in internal/schemas/manager_test.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/manager_test.go
```

### internal/schemas/registry.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewValidator

```bash
# TODO: Fix linting issues in internal/schemas/registry.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/registry.go
```

### internal/secrets/integration.go

**Issues found:** 3

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/secrets/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/secrets/integration.go
```

### internal/ssh/connection_pool.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: PooledConnection
- **typecheck** (High): undefined: PooledConnection
- **typecheck** (High): undefined: PooledConnection

```bash
# TODO: Fix linting issues in internal/ssh/connection_pool.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/connection_pool.go
```

### internal/ssh/file_transfer.go

**Issues found:** 4

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Client
- **typecheck** (High): undefined: Client
- **typecheck** (High): undefined: PooledConnection
- **typecheck** (High): undefined: PooledConnection

```bash
# TODO: Fix linting issues in internal/ssh/file_transfer.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/file_transfer.go
```

### internal/ssh/host_key_test.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewHostKeyManager
- **typecheck** (High): undefined: NewClient
- **typecheck** (High): undefined: NewClient

```bash
# TODO: Fix linting issues in internal/ssh/host_key_test.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/host_key_test.go
```

### internal/ssh/manager_test.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewManager
- **typecheck** (High): undefined: NewManager

```bash
# TODO: Fix linting issues in internal/ssh/manager_test.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/manager_test.go
```

### internal/templates/integration.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Manager
- **typecheck** (High): undefined: Manager

```bash
# TODO: Fix linting issues in internal/templates/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/templates/integration.go
```

### internal/templates/manager.go

**Issues found:** 2

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/templates/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/templates/manager.go
```

### internal/types/actions/structures.go

**Issues found:** 4

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Action
- **typecheck** (High): undefined: Action
- **typecheck** (High): undefined: Action
- **typecheck** (High): undefined: ValidationError

```bash
# TODO: Fix linting issues in internal/types/actions/structures.go
# Run: golangci-lint run --config=.golangci.yml internal/types/actions/structures.go
```

### internal/types/config/config.go

**Issues found:** 4

#### 🟢 Low Priority Issues

- **revive** (Low): exported: type name will be used as config.ConfigDefaults by other packages, and that stutters; consider calling this Defaults
- **revive** (Low): exported: type name will be used as config.ConfigSupporting by other packages, and that stutters; consider calling this Supporting
- **revive** (Low): exported: type name will be used as config.ConfigEnvironment by other packages, and that stutters; consider calling this Environment
- **revive** (Low): exported: type name will be used as config.ConfigError by other packages, and that stutters; consider calling this Error

```bash
# TODO: Fix linting issues in internal/types/config/config.go
# Run: golangci-lint run --config=.golangci.yml internal/types/config/config.go
```

### internal/types/facts/collection.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Facts

```bash
# TODO: Fix linting issues in internal/types/facts/collection.go
# Run: golangci-lint run --config=.golangci.yml internal/types/facts/collection.go
```

### internal/types/ssh/authentication.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: AuthMethod
- **typecheck** (High): undefined: AuthMethod

```bash
# TODO: Fix linting issues in internal/types/ssh/authentication.go
# Run: golangci-lint run --config=.golangci.yml internal/types/ssh/authentication.go
```

### internal/types/ssh/client.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: AuthMethod
- **typecheck** (High): undefined: ConnectionPool
- **typecheck** (High): undefined: Connection

```bash
# TODO: Fix linting issues in internal/types/ssh/client.go
# Run: golangci-lint run --config=.golangci.yml internal/types/ssh/client.go
```

### internal/variables/integration.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewManager

```bash
# TODO: Fix linting issues in internal/variables/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/integration.go
```

### internal/variables/loader.go

**Issues found:** 2

#### 🟢 Low Priority Issues

- **gocyclo** (Low): cyclomatic complexity 23 of func `(*Loader).parseVariableBlock` is high (> 20)
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/variables/loader.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/loader.go
```

### internal/variables/manager.go

**Issues found:** 2

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): unused-parameter: parameter 'context' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/variables/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/manager.go
```

### internal/variables/validator.go

**Issues found:** 4

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): unnamedResult: consider giving a name to these results

#### 🟢 Low Priority Issues

- **gocyclo** (Low): cyclomatic complexity 27 of func `(*Validator).validateVariableConstraints` is high (> 20)
- **revive** (Low): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (Low): redefines-builtin-id: redefinition of the built-in type error

```bash
# TODO: Fix linting issues in internal/variables/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/validator.go
```

### scripts/lint-todos.go

**Issues found:** 4

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): nestingReduce: invert if cond, replace body with `continue`, move old body after the statement

#### 🟢 Low Priority Issues

- **revive** (Low): unused-parameter: parameter 'message' seems to be unused, consider removing or renaming it as _
- **revive** (Low): superfluous-else: if block ends with call to os.Exit function, so drop this else and outdent its block
- **revive** (Low): superfluous-else: if block ends with call to os.Exit function, so drop this else and outdent its block

```bash
# TODO: Fix linting issues in scripts/lint-todos.go
# Run: golangci-lint run --config=.golangci.yml scripts/lint-todos.go
```

## 🔧 Batch Fix Opportunities

### Issues Found in Multiple Files

#### revive (17 files)

**Pattern:** unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

**Affected files:**
- internal/cli/commands/integrations.go:25
- internal/machines/validator.go:62
- internal/project/loader.go:79
- internal/project/loader.go:374
- internal/project/manager.go:39
- internal/project/manager.go:173
- internal/project/manager.go:192
- internal/project/validator.go:83
- internal/project/validator.go:185
- internal/secrets/integration.go:29
- internal/secrets/integration.go:68
- internal/secrets/integration.go:112
- internal/templates/manager.go:65
- internal/templates/manager.go:102
- internal/variables/loader.go:32
- internal/variables/manager.go:135
- internal/variables/validator.go:66

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (11 files)

**Pattern:** undefined: NewManager

**Affected files:**
- internal/facts/integration_test.go:19
- internal/facts/integration_test.go:37
- internal/facts/integration_test.go:71
- internal/facts/manager_test.go:94
- internal/schemas/manager_test.go:17
- internal/schemas/manager_test.go:58
- internal/schemas/manager_test.go:91
- internal/ssh/manager_test.go:18
- internal/ssh/manager_test.go:30
- internal/ssh/manager_test.go:79
- internal/variables/integration.go:22

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (9 files)

**Pattern:** rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)

**Affected files:**
- internal/actions/manager_test.go:390
- internal/actions/manager_test.go:402
- internal/machines/integration.go:107
- internal/machines/integration.go:218
- internal/machines/integration.go:250
- internal/machines/manager.go:131
- internal/machines/manager.go:391
- internal/machines/manager.go:551
- internal/machines/validator.go:39

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (9 files)

**Pattern:** undefined: FactCollection

**Affected files:**
- internal/facts/collector.go:38
- internal/facts/collector.go:76
- internal/facts/integration.go:58
- internal/facts/integration.go:96
- internal/facts/integration_test.go:74
- internal/facts/integration_test.go:122
- internal/facts/manager_test.go:16
- internal/facts/manager_test.go:25
- internal/facts/manager_test.go:30

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (5 files)

**Pattern:** undefined: Manager

**Affected files:**
- internal/facts/integration.go:16
- internal/facts/integration.go:21
- internal/facts/manager_test.go:89
- internal/templates/integration.go:16
- internal/templates/integration.go:20

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (5 files)

**Pattern:** undefined: PooledConnection

**Affected files:**
- internal/ssh/connection_pool.go:17
- internal/ssh/connection_pool.go:52
- internal/ssh/connection_pool.go:92
- internal/ssh/file_transfer.go:102
- internal/ssh/file_transfer.go:285

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (4 files)

**Pattern:** hugeParam: machine is heavy (512 bytes); consider passing it by pointer

**Affected files:**
- internal/machines/integration.go:382
- internal/machines/manager.go:436
- internal/machines/validator.go:175
- internal/machines/validator.go:197

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (4 files)

**Pattern:** undefined: MockLogger

**Affected files:**
- internal/facts/integration_test.go:17
- internal/facts/integration_test.go:35
- internal/facts/integration_test.go:69
- internal/facts/manager_test.go:92

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (4 files)

**Pattern:** undefined: MockSchemaValidator

**Affected files:**
- internal/facts/integration_test.go:16
- internal/facts/integration_test.go:34
- internal/facts/integration_test.go:68
- internal/facts/manager_test.go:91

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (4 files)

**Pattern:** octalLiteral: use new octal literal style, 0o755

**Affected files:**
- internal/machines/manager.go:250
- internal/project/manager.go:56
- internal/project/manager.go:66
- internal/project/manager.go:78

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: AuthMethod

**Affected files:**
- internal/types/ssh/authentication.go:16
- internal/types/ssh/authentication.go:210
- internal/types/ssh/client.go:47

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (3 files)

**Pattern:** sprintfQuotedString: use %q instead of "%s" for quoted strings

**Affected files:**
- internal/machines/integration.go:341
- internal/machines/integration.go:342
- internal/machines/integration.go:344

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: NewIntegration

**Affected files:**
- internal/facts/integration_test.go:21
- internal/facts/integration_test.go:38
- internal/facts/integration_test.go:72

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: Action

**Affected files:**
- internal/types/actions/structures.go:15
- internal/types/actions/structures.go:36
- internal/types/actions/structures.go:69

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: MockFactCollector

**Affected files:**
- internal/facts/integration_test.go:15
- internal/facts/integration_test.go:33
- internal/facts/integration_test.go:67

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

## 📈 Progress Tracking

### Current Status

- **Fixed:** 0/139 issues (0%)
- **High Priority:** 69 issues (49.6%)
- **Medium Priority:** 31 issues (22.3%)
- **Low Priority:** 39 issues (28.1%)

### Recommended Action Plan

1. **🔴 Fix all High Priority issues first** (69 issues) - prevents compilation
2. **🟡 Address Medium Priority issues** (31 issues) - improves code quality
3. **🟢 Fix Low Priority issues** (39 issues) - style and formatting
4. **🔧 Apply batch fixes** where possible
5. **✅ Run `golangci-lint run` to verify all fixes**

## 🛠️ Common Fix Patterns

### For "undefined" errors:
```go
// Pattern 1: Add missing import
import "spooky/internal/missingpackage"

// Pattern 2: Define missing variable
var missingVar = NewMissingType()

// Pattern 3: Fix function call
existingFunction() // instead of undefinedFunction()
```

### For "hugeParam" warnings:
```go
// Before: func process(config Config)
// After: func process(config *Config)
```

### For "unused" warnings:
```go
// Remove unused imports
// import _ "unusedpackage" // Remove this line

// Remove unused variables
// var unusedVar = "value" // Remove this line
```

### For formatting issues:
```bash
# Run gofmt to fix formatting
gofmt -w file.go

# Run goimports to fix imports
goimports -w file.go
```

