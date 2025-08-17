# Dependency Analysis

Generated on: Sun 17 Aug 2025 08:06:48 CEST

## Summary

- **Total packages analyzed:** 14
- **Total interfaces:** 0
- **Total types:** 0
- **Total functions:** 0

## Generated Reports

- [Package Dependencies](package-dependencies.md) - Overall package dependency graph
- [Interface Dependencies](interface-dependencies.md) - Interface relationships
- [Type Dependencies](type-dependencies.md) - Type relationships
- [Function Dependencies](function-dependencies.md) - Function relationships

## Usage

```bash
# Generate all dependency graphs
just dependency-graph

# Analyze specific package
just dependency-graph-package internal/ssh

# Generate specific graph type
just dependency-graph-interface
just dependency-graph-type
just dependency-graph-function
```
