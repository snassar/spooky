# Implementation Plan: Types Consolidation

## Overview

Consolidate all package-specific types into a unified `./internal/types/` package structure to resolve import cycles, improve consistency, and simplify the codebase architecture.

## Current State

### ✅ **COMPLETED: Types Consolidation**
- **Unified types package**: `internal/types/` structure created and populated
- **All package-specific types**: Moved to `internal/types/<package>/` subdirectories
- **Import aliases updated**: Using unified `spookytypes "spooky/internal/types"` pattern
- **Documentation updated**: Import alias enforcement rule reflects current patterns

### ✅ **COMPLETED: Import Pattern Migration**
- **Old pattern deprecated**: `spooky<package>types` (e.g., `spookyconfigtypes`, `spookyfactstypes`)
- **New pattern active**: `spookytypes "spooky/internal/types"` for all types
- **Consistent usage**: All internal packages use unified types package

## Original Problems (RESOLVED)

### 1. **Import Cycle Issues** ✅ **RESOLVED**
- ~~`internal/config/loading/parser.go` imports `spooky/internal/schemas`~~
- ~~`internal/schemas/project.go` imports `spooky/internal/project`~~
- ~~`internal/project/configuration.go` imports `spooky/internal/config`~~
- **Result**: ✅ Circular dependencies resolved through unified types package

### 2. **Inconsistent Patterns** ✅ **RESOLVED**
- ~~Some packages have `types/` subdirectories, others don't~~
- ~~Types are scattered across multiple packages~~
- ~~Import aliases are complex and inconsistent~~
- **Result**: ✅ All types now in unified `internal/types/` structure

### 3. **Architectural Issues** ✅ **RESOLVED**
- ~~Package-specific types create tight coupling~~
- ~~Difficult to share types between packages~~
- ~~Import complexity with long aliases~~
- **Result**: ✅ Single source of truth for all types with simplified aliases

## Implemented Solution

### **Current Structure:**
```
internal/
├── types/                           # ✅ Unified types package
│   ├── common/                      # ✅ Shared common types
│   │   └── common.go               # EncryptionMetadata, TimestampedEntity, etc.
│   ├── actions/                     # ✅ Actions-specific types
│   │   ├── action.go
│   │   ├── context.go
│   │   └── collection_planning.go
│   ├── config/                      # ✅ Config-specific types
│   │   ├── config.go
│   │   ├── defaults.go
│   │   └── environment.go
│   ├── facts/                       # ✅ Facts-specific types
│   │   ├── facts.go
│   │   ├── collection.go
│   │   └── errors.go
│   ├── machines/                    # ✅ Machines-specific types
│   │   ├── config.go
│   │   ├── errors.go
│   │   └── import_export.go
│   ├── ssh/                         # ✅ SSH-specific types
│   │   ├── acting.go
│   │   ├── authentication.go
│   │   └── client.go
│   ├── templates/                   # ✅ Templates-specific types
│   │   ├── context.go
│   │   ├── errors.go
│   │   ├── functions.go
│   │   └── template.go
│   ├── variables/                   # ✅ Variables-specific types
│   │   └── variable.go
│   ├── logging/                     # ✅ Logging-specific types
│   │   ├── config.go
│   │   ├── entry.go
│   │   └── errors.go
│   ├── secrets/                     # ✅ Secrets-specific types
│   │   └── types.go
│   ├── schemas/                     # ✅ Schemas-specific types
│   │   ├── errors.go
│   │   ├── schema.go
│   │   └── validation.go
│   ├── cli/                         # ✅ CLI-specific types
│   │   ├── command.go
│   │   ├── context.go
│   │   └── errors.go
│   └── project/                     # ✅ Project-specific types
│       └── structures.go
├── actions/                         # ✅ Types moved to internal/types/actions/
├── config/                          # ✅ Types moved to internal/types/config/
├── facts/                           # ✅ Types moved to internal/types/facts/
├── machines/                        # ✅ Types moved to internal/types/machines/
├── ssh/                             # ✅ Types moved to internal/types/ssh/
├── templates/                       # ✅ Types moved to internal/types/templates/
├── variables/                       # ✅ Types moved to internal/types/variables/
├── logging/                         # ✅ Types moved to internal/types/logging/
├── secrets/                         # ✅ Types moved to internal/types/secrets/
├── schemas/                         # ✅ Types moved to internal/types/schemas/
├── cli/                             # ✅ Types moved to internal/types/cli/
└── project/                         # ✅ Types moved to internal/types/project/
```

## Implementation Status

### ✅ **Phase 1: Create New Types Structure** - **COMPLETED**
1. ✅ Create `internal/types/` directory structure
2. ✅ Move `internal/types/common/` (already exists)
3. ✅ Create package-specific subdirectories
4. ✅ Move types from each package to `internal/types/<package>/`

### ✅ **Phase 2: Update Imports** - **COMPLETED**
1. ✅ Update all import statements to use new types paths
2. ✅ Simplify import aliases to `spookytypes "spooky/internal/types"`
3. ✅ Remove old types directories
4. ✅ Update documentation

### ✅ **Phase 3: Fix Import Cycles** - **COMPLETED**
1. ✅ Resolve circular dependencies
2. ✅ Enable schema validation in config parser
3. ✅ Test all packages build successfully

### ✅ **Phase 4: Update Documentation** - **COMPLETED**
1. ✅ Update import alias documentation
2. ✅ Update package architecture rules
3. ✅ Update examples and tests

## Current Import Patterns

### ✅ **Unified Types Pattern (CURRENT)**
```go
// ✅ CORRECT - Unified types package
import (
    spookytypes "spooky/internal/types"
)

// Usage examples
func (c *Config) Validate() error {
    return spookytypes.ValidateConfig(c)
}

func NewManager() *spookytypes.Manager {
    return spookytypes.NewManager()
}

func ConfigureLogger(config *spookytypes.LoggingConfig) {
    // ...
}
```

### ❌ **Old Pattern (DEPRECATED)**
```go
// ❌ DEPRECATED - Old pattern (DO NOT USE)
import (
    spookyconfigtypes "spooky/internal/config/types"
    spookyfactstypes "spooky/internal/facts/types"
    spookyactionstypes "spooky/internal/actions/types"
)
```

## Benefits Achieved

### ✅ **1. Resolves Import Cycles**
- Types are now independent of package implementations
- No more circular dependencies between packages
- Schema validation can be used everywhere

### ✅ **2. Simplifies Architecture**
- Single source of truth for all types
- Consistent import patterns with `spookytypes` alias
- Easier to understand and maintain

### ✅ **3. Improves Developer Experience**
- Shorter, cleaner import aliases (`spookytypes` vs `spooky<package>types`)
- Better IDE support and autocomplete
- Reduced cognitive load

### ✅ **4. Enables Future Features**
- Schema validation everywhere
- Better type sharing between packages
- Cleaner dependency graph

## Migration Completed

### ✅ **Step 1: Create New Structure** - **COMPLETED**
```bash
# ✅ Completed: Create new types directory structure
mkdir -p internal/types/{actions,config,facts,machines,ssh,templates,variables,logging,secrets,schemas,cli,project}

# ✅ Completed: Move existing common types
# internal/types/common/common.go
```

### ✅ **Step 2: Move Types by Package** - **COMPLETED**
```bash
# ✅ Completed: Move config types
mv internal/config/types/* internal/types/config/
rmdir internal/config/types

# ✅ Completed: Repeat for all packages
```

### ✅ **Step 3: Update Imports** - **COMPLETED**
```go
// ✅ COMPLETED: Before
import (
    spookyconfigtypes "spooky/internal/config/types"
    spookyfactstypes "spooky/internal/facts/types"
)

// ✅ COMPLETED: After
import (
    spookytypes "spooky/internal/types"
)

// ✅ COMPLETED: Usage
config := &spookytypes.Config{}
facts := &spookytypes.FactCollection{}
```

### ✅ **Step 4: Update Documentation** - **COMPLETED**
- ✅ Update `docs/import-aliases.md`
- ✅ Update `.cursor/rules/package-architecture.mdc`
- ✅ Update `.cursor/rules/import-alias-enforcement.mdc`
- ✅ Update all examples and tests

## Success Criteria - ALL MET ✅

1. ✅ **No Import Cycles**: All packages build without circular dependencies
2. ✅ **Schema Validation Works**: Config parser can use schema validation
3. ✅ **Cleaner Imports**: Simplified import aliases throughout codebase (`spookytypes`)
4. ✅ **Consistent Structure**: All types follow same organizational pattern
5. ✅ **Documentation Updated**: All docs reflect new structure
6. ✅ **Tests Pass**: All existing tests continue to work

## Enforcement Rules

### **Import Alias Enforcement**
- **Rule**: `.cursor/rules/import-alias-enforcement.mdc`
- **Pattern**: `spookytypes "spooky/internal/types"` for all types
- **Status**: ✅ Active and enforced

### **Forbidden Patterns**
```go
// ❌ DO NOT USE - Old patterns
import (
    spookyconfigtypes "spooky/internal/config/types"    // ❌ Deprecated
    spookyfactstypes "spooky/internal/facts/types"      // ❌ Deprecated
    spookyactionstypes "spooky/internal/actions/types"  // ❌ Deprecated
)

// ❌ DO NOT USE - Naked imports
import (
    "spooky/internal/config/types"    // ❌ Naked import
    "spooky/internal/facts/types"     // ❌ Naked import
)
```

### **Required Patterns**
```go
// ✅ USE THIS - Unified types pattern
import (
    spookytypes "spooky/internal/types"
)
```

## Timeline - COMPLETED ✅

- ✅ **Phase 1**: 1-2 days (create structure, move types) - **COMPLETED**
- ✅ **Phase 2**: 1-2 days (update imports, remove old directories) - **COMPLETED**
- ✅ **Phase 3**: 1 day (fix cycles, test builds) - **COMPLETED**
- ✅ **Phase 4**: 1 day (update docs, final testing) - **COMPLETED**

**Total**: ✅ **COMPLETED** - Types consolidation successfully implemented

## Next Steps

### **Maintenance**
1. **Enforce unified pattern**: Use `spookytypes` alias for all new code
2. **Monitor for regressions**: Ensure no old patterns are reintroduced
3. **Update documentation**: Keep docs current with any new types

### **Future Enhancements**
1. **Type validation**: Add validation rules for new types
2. **Type documentation**: Improve type documentation and examples
3. **Type testing**: Add comprehensive tests for all types

---

**Status**: ✅ **COMPLETED** - Types consolidation successfully implemented and enforced
