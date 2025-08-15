# Linting TODOs - AI-Optimized Report

This file contains linting issues captured by running `golangci-lint` on each Go file individually.
**Optimized for AI consumption** with categorization, fix suggestions, and progress tracking.

**Generated:** 2025-08-15 07:51:36
**Config:** .golangci.yml
**Fast Mode:** true

## 📊 Summary Statistics

- **Total files processed:** 88
- **Files with issues:** 35
- **Total issues:** 124

### Priority Breakdown

- **🔴 High Priority (Build-breaking):** 84 issues (67.7%)
- **🟡 Medium Priority (Code quality):** 40 issues (32.3%)
- **🟢 Low Priority (Style):** 0 issues (0.0%)

## 🏷️ Issue Categories

### typecheck (61 issues)

**Description:** Build-breaking errors that prevent compilation

**Common Fix Patterns:**
- Add missing import: `import "package/path"`
- Define missing variable: `var variableName Type`
- Fix function call: use correct function name
- Add missing package initialization

**Examples:**
- `cmd/actions.go:277:9`: undefined: integrationManager
- `cmd/actions.go:281:2`: undefined: RootCmd
- `cmd/integrations.go:100:2`: undefined: RootCmd

### gocritic (27 issues)

**Description:** Code quality and style suggestions

**Common Fix Patterns:**
- Use pointer for large structs: `func process(config *Config)`
- Simplify boolean expressions
- Use consistent naming conventions
- Optimize string concatenation

**Examples:**
- `internal/actions/manager_test.go:390:2`: rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- `internal/actions/manager_test.go:402:2`: rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- `internal/cli/commands/integrations.go:79:2`: rangeValCopy: each iteration copies 248 bytes (consider pointers or indexing)

### revive (26 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/cli/commands/integrations.go:25:36`: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- `internal/cli/commands/integrations.go:57:61`: unused-parameter: parameter 'args' seems to be unused, consider removing or renaming it as _
- `internal/logging/logging.go:463:43`: unused-parameter: parameter 'fields' seems to be unused, consider removing or renaming it as _

### gocyclo (5 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/facts/hcl_parser.go:143:1`: cyclomatic complexity 24 of func `(*HCLParser).parseHostBlock` is high (> 20)
- `internal/machines/manager.go:116:1`: cyclomatic complexity 34 of func `(*Manager).ExportMachines` is high (> 20)
- `internal/machines/manager.go:804:1`: cyclomatic complexity 27 of func `(*Manager).applyMachineFilter` is high (> 20)

### gosimple (2 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/machines/manager.go:770:5`: S1009: should omit nil check; len() for nil maps is defined as zero
- `internal/schemas/enhanced_validator.go:402:5`: S1009: should omit nil check; len() for nil slices is defined as zero

### ineffassign (1 issues)

**Description:** Ineffective assignments

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/logging/logging.go:484:6`: ineffectual assignment to matched

### gosec (1 issues)

**Description:** Code quality issue

**Common Fix Patterns:**
- Review and fix according to linter suggestions

**Examples:**
- `internal/machines/manager.go:255:12`: G306: Expect WriteFile permissions to be 0600 or less

### unused (1 issues)

**Description:** Unused variables, imports, or functions

**Common Fix Patterns:**
- Remove unused imports: `import _ "package"`
- Remove unused variables
- Remove unused functions
- Use variables or remove them

**Examples:**
- `internal/machines/manager.go:549:6`: func `getMachineHostnames` is unused

## 📁 File-by-File Breakdown

### cmd/actions.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: integrationManager
- **typecheck** (High): undefined: RootCmd

```bash
# TODO: Fix linting issues in cmd/actions.go
# Run: golangci-lint run --config=.golangci.yml cmd/actions.go
```

### cmd/integrations.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: RootCmd

```bash
# TODO: Fix linting issues in cmd/integrations.go
# Run: golangci-lint run --config=.golangci.yml cmd/integrations.go
```

### cmd/project.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: RootCmd

```bash
# TODO: Fix linting issues in cmd/project.go
# Run: golangci-lint run --config=.golangci.yml cmd/project.go
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

### internal/cli/commands/integrations.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'args' seems to be unused, consider removing or renaming it as _

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): rangeValCopy: each iteration copies 248 bytes (consider pointers or indexing)

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

#### 🟡 Medium Priority Issues

- **gocyclo** (Medium): cyclomatic complexity 24 of func `(*HCLParser).parseHostBlock` is high (> 20)

```bash
# TODO: Fix linting issues in internal/facts/hcl_parser.go
# Run: golangci-lint run --config=.golangci.yml internal/facts/hcl_parser.go
```

### internal/logging/logging.go

**Issues found:** 12

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'fields' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'attrs' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'name' seems to be unused, consider removing or renaming it as _

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): appendCombine: can combine chain of 3 appends into one
- **gocritic** (Medium): hugeParam: r is heavy (288 bytes); consider passing it by pointer
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o644
- **gocritic** (Medium): sprintfQuotedString: use %q instead of "%s" for quoted strings
- **gocritic** (Medium): sprintfQuotedString: use %q instead of "%s" for quoted strings
- **ineffassign** (Medium): ineffectual assignment to matched

```bash
# TODO: Fix linting issues in internal/logging/logging.go
# Run: golangci-lint run --config=.golangci.yml internal/logging/logging.go
```

### internal/logging/logging_test.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewLogManager
- **typecheck** (High): undefined: NewLogManager
- **typecheck** (High): undefined: NewLogManager

```bash
# TODO: Fix linting issues in internal/logging/logging_test.go
# Run: golangci-lint run --config=.golangci.yml internal/logging/logging_test.go
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

#### 🔴 High Priority Issues

- **unused** (High): func `getMachineHostnames` is unused

#### 🟡 Medium Priority Issues

- **gocyclo** (Medium): cyclomatic complexity 34 of func `(*Manager).ExportMachines` is high (> 20)
- **gocyclo** (Medium): cyclomatic complexity 27 of func `(*Manager).applyMachineFilter` is high (> 20)
- **gocyclo** (Medium): cyclomatic complexity 22 of func `(*Manager).validateMachineCollection` is high (> 20)
- **gosec** (Medium): G306: Expect WriteFile permissions to be 0600 or less
- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): ifElseChain: rewrite if-else to switch statement
- **gocritic** (Medium): nestingReduce: invert if cond, replace body with `continue`, move old body after the statement
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gosimple** (Medium): S1009: should omit nil check; len() for nil maps is defined as zero

```bash
# TODO: Fix linting issues in internal/machines/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/manager.go
```

### internal/machines/validator.go

**Issues found:** 5

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): hugeParam: machine is heavy (512 bytes); consider passing it by pointer
- **gocritic** (Medium): rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)
- **gocritic** (Medium): regexpSimplify: can re-write `^[0-9]+(Gbps|Mbps)$` as `^\d+(Gbps|Mbps)$`

```bash
# TODO: Fix linting issues in internal/machines/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/machines/validator.go
```

### internal/project/loader.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'projectPath' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/project/loader.go
# Run: golangci-lint run --config=.golangci.yml internal/project/loader.go
```

### internal/project/manager.go

**Issues found:** 6

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

#### 🟡 Medium Priority Issues

- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755
- **gocritic** (Medium): octalLiteral: use new octal literal style, 0o755

```bash
# TODO: Fix linting issues in internal/project/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/project/manager.go
```

### internal/project/validator.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/project/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/project/validator.go
```

### internal/schemas/enhanced_validator.go

**Issues found:** 9

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'schema' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'schema' seems to be unused, consider removing or renaming it as _

#### 🟡 Medium Priority Issues

- **gocyclo** (Medium): cyclomatic complexity 21 of func `(*EnhancedValidator).validateFieldConstraintsValue` is high (> 20)
- **revive** (Medium): redefines-builtin-id: redefinition of the built-in type error
- **revive** (Medium): redefines-builtin-id: redefinition of the built-in type error
- **revive** (Medium): redefines-builtin-id: redefinition of the built-in type error
- **gocritic** (Medium): rangeValCopy: each iteration copies 248 bytes (consider pointers or indexing)
- **gocritic** (Medium): typeAssertChain: rewrite if-else to type switch statement
- **gosimple** (Medium): S1009: should omit nil check; len() for nil slices is defined as zero

```bash
# TODO: Fix linting issues in internal/schemas/enhanced_validator.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/enhanced_validator.go
```

### internal/schemas/evolution_manager.go

**Issues found:** 3

#### 🟡 Medium Priority Issues

- **revive** (Medium): redefines-builtin-id: redefinition of the built-in type error
- **gocritic** (Medium): rangeValCopy: each iteration copies 136 bytes (consider pointers or indexing)
- **gocritic** (Medium): rangeValCopy: each iteration copies 136 bytes (consider pointers or indexing)

```bash
# TODO: Fix linting issues in internal/schemas/evolution_manager.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/evolution_manager.go
```

### internal/schemas/manager.go

**Issues found:** 8

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Registry
- **typecheck** (High): undefined: Validator
- **typecheck** (High): undefined: EnhancedValidator
- **typecheck** (High): undefined: EvolutionManager
- **typecheck** (High): undefined: NewRegistry
- **typecheck** (High): undefined: NewValidator
- **typecheck** (High): undefined: NewEnhancedValidator
- **typecheck** (High): undefined: NewEvolutionManager

```bash
# TODO: Fix linting issues in internal/schemas/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/manager.go
```

### internal/schemas/registry.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewValidator

```bash
# TODO: Fix linting issues in internal/schemas/registry.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/registry.go
```

### internal/schemas/validator.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Manager
- **typecheck** (High): undefined: Manager

```bash
# TODO: Fix linting issues in internal/schemas/validator.go
# Run: golangci-lint run --config=.golangci.yml internal/schemas/validator.go
```

### internal/ssh/connection_pool_test.go

**Issues found:** 6

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewAdvancedConnectionPool
- **typecheck** (High): undefined: NewAdvancedConnectionPool
- **typecheck** (High): undefined: NewAdvancedConnectionPool
- **typecheck** (High): undefined: PooledConnection
- **typecheck** (High): undefined: PooledConnection
- **typecheck** (High): undefined: PooledConnection

```bash
# TODO: Fix linting issues in internal/ssh/connection_pool_test.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/connection_pool_test.go
```

### internal/ssh/example_test.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewClient

```bash
# TODO: Fix linting issues in internal/ssh/example_test.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/example_test.go
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

### internal/ssh/manager.go

**Issues found:** 3

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Client
- **typecheck** (High): undefined: NewClient
- **typecheck** (High): undefined: NewClient

```bash
# TODO: Fix linting issues in internal/ssh/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/ssh/manager.go
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

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/templates/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/templates/manager.go
```

### internal/types/ssh/acting.go

**Issues found:** 5

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: Connection
- **typecheck** (High): undefined: Client
- **typecheck** (High): undefined: PtyConfig
- **typecheck** (High): undefined: PtyConfig
- **typecheck** (High): undefined: PtyConfig

```bash
# TODO: Fix linting issues in internal/types/ssh/acting.go
# Run: golangci-lint run --config=.golangci.yml internal/types/ssh/acting.go
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

### internal/types/ssh/errors.go

**Issues found:** 5

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: AuthMethod
- **typecheck** (High): undefined: KeyType
- **typecheck** (High): undefined: KeyType
- **typecheck** (High): undefined: TransferMode
- **typecheck** (High): undefined: TransferDirection

```bash
# TODO: Fix linting issues in internal/types/ssh/errors.go
# Run: golangci-lint run --config=.golangci.yml internal/types/ssh/errors.go
```

### internal/variables/integration.go

**Issues found:** 1

#### 🔴 High Priority Issues

- **typecheck** (High): undefined: NewManager

```bash
# TODO: Fix linting issues in internal/variables/integration.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/integration.go
```

### internal/variables/manager.go

**Issues found:** 2

#### 🔴 High Priority Issues

- **revive** (High): unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _
- **revive** (High): unused-parameter: parameter 'context' seems to be unused, consider removing or renaming it as _

```bash
# TODO: Fix linting issues in internal/variables/manager.go
# Run: golangci-lint run --config=.golangci.yml internal/variables/manager.go
```

## 🔧 Batch Fix Opportunities

### Issues Found in Multiple Files

#### revive (14 files)

**Pattern:** unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _

**Affected files:**
- internal/cli/commands/integrations.go:25
- internal/logging/logging.go:547
- internal/logging/logging.go:552
- internal/machines/validator.go:62
- internal/project/loader.go:79
- internal/project/loader.go:374
- internal/project/manager.go:39
- internal/project/manager.go:173
- internal/project/manager.go:192
- internal/project/validator.go:83
- internal/project/validator.go:185
- internal/templates/manager.go:65
- internal/templates/manager.go:102
- internal/variables/manager.go:135

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (6 files)

**Pattern:** rangeValCopy: each iteration copies 512 bytes (consider pointers or indexing)

**Affected files:**
- internal/actions/manager_test.go:390
- internal/actions/manager_test.go:402
- internal/machines/manager.go:131
- internal/machines/manager.go:391
- internal/machines/manager.go:551
- internal/machines/validator.go:39

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (5 files)

**Pattern:** octalLiteral: use new octal literal style, 0o755

**Affected files:**
- internal/logging/logging.go:168
- internal/machines/manager.go:250
- internal/project/manager.go:56
- internal/project/manager.go:66
- internal/project/manager.go:78

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (5 files)

**Pattern:** undefined: NewClient

**Affected files:**
- internal/ssh/example_test.go:32
- internal/ssh/host_key_test.go:74
- internal/ssh/host_key_test.go:111
- internal/ssh/manager.go:43
- internal/ssh/manager.go:56

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### revive (4 files)

**Pattern:** redefines-builtin-id: redefinition of the built-in type error

**Affected files:**
- internal/schemas/enhanced_validator.go:456
- internal/schemas/enhanced_validator.go:494
- internal/schemas/enhanced_validator.go:599
- internal/schemas/evolution_manager.go:173

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (4 files)

**Pattern:** undefined: Manager

**Affected files:**
- internal/schemas/validator.go:25
- internal/schemas/validator.go:42
- internal/templates/integration.go:16
- internal/templates/integration.go:20

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (4 files)

**Pattern:** undefined: NewManager

**Affected files:**
- internal/ssh/manager_test.go:18
- internal/ssh/manager_test.go:30
- internal/ssh/manager_test.go:79
- internal/variables/integration.go:22

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (4 files)

**Pattern:** undefined: AuthMethod

**Affected files:**
- internal/types/ssh/authentication.go:16
- internal/types/ssh/authentication.go:210
- internal/types/ssh/client.go:47
- internal/types/ssh/errors.go:79

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: NewLogManager

**Affected files:**
- internal/logging/logging_test.go:17
- internal/logging/logging_test.go:54
- internal/logging/logging_test.go:87

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: PtyConfig

**Affected files:**
- internal/types/ssh/acting.go:28
- internal/types/ssh/acting.go:80
- internal/types/ssh/acting.go:231

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: PooledConnection

**Affected files:**
- internal/ssh/connection_pool_test.go:158
- internal/ssh/connection_pool_test.go:180
- internal/ssh/connection_pool_test.go:253

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: NewAdvancedConnectionPool

**Affected files:**
- internal/ssh/connection_pool_test.go:30
- internal/ssh/connection_pool_test.go:70
- internal/ssh/connection_pool_test.go:111

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### gocritic (3 files)

**Pattern:** hugeParam: machine is heavy (512 bytes); consider passing it by pointer

**Affected files:**
- internal/machines/manager.go:436
- internal/machines/validator.go:175
- internal/machines/validator.go:197

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

#### typecheck (3 files)

**Pattern:** undefined: RootCmd

**Affected files:**
- cmd/actions.go:281
- cmd/integrations.go:100
- cmd/project.go:198

**Suggested batch fix:**
```bash
# Run this command to fix all instances:
# find . -name "*.go" -exec sed -i 's/pattern/replacement/g' {} \;
```

## 📈 Progress Tracking

### Current Status

- **Fixed:** 0/124 issues (0%)
- **High Priority:** 84 issues (67.7%)
- **Medium Priority:** 40 issues (32.3%)
- **Low Priority:** 0 issues (0.0%)

### Recommended Action Plan

1. **🔴 Fix all High Priority issues first** (84 issues) - prevents compilation
2. **🟡 Address Medium Priority issues** (40 issues) - improves code quality
3. **🟢 Fix Low Priority issues** (0 issues) - style and formatting
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

