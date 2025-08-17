# Architecture Decision Records (ADRs)

This directory contains Architecture Decision Records (ADRs) for the spooky project. ADRs document significant architectural decisions, their context, and consequences.

## What are ADRs?

Architecture Decision Records are documents that capture important architectural decisions made during the development of a project. They include:

- **Context**: The problem being solved
- **Decision**: What was chosen and why
- **Consequences**: Trade-offs and implications
- **Evidence**: Supporting information from git history

## Generated ADRs

The ADRs in `docs/adr/data/` are automatically generated from git history analysis. They reconstruct architectural decisions based on:

- Commit messages and patterns
- File structure evolution
- Code changes and refactoring
- Interface and type system development

## Analysis Tools

### Focused Analysis (`docs/adr/recommendations/`)
Identifies only the most significant architectural decisions that should be documented:

- **High Priority**: Breaking changes, major refactors, architecture changes
- **Medium Priority**: System changes, interface changes, large refactors
- **Actionable**: Only recommendations that would benefit from ADR documentation

### Comprehensive Analysis (`docs/adr/analysis/`)
Provides a complete view of all potential architectural decisions:

- **All commits** with architectural keywords
- **Confidence scoring** for each potential ADR
- **Categorized by type** (Architecture, Interface, System, etc.)
- **Useful for understanding** project evolution and decision patterns

### Simple ADRs (`docs/adr/`)
Direct ADR generation with full paths from spooky root:

- **Simple format**: Clean ADR files directly in `docs/adr/`
- **Full paths**: Uses complete paths from spooky root
- **One command**: `just generate-adrs-simple`
- **Generated from git**: Automatically identifies significant commits

## Categories

ADRs are organized into categories:

- **Architecture**: Core architectural patterns and decisions
- **Features**: Major feature implementations
- **Refactoring**: Structural changes and improvements
- **Interfaces**: Interface design and coordination patterns
- **Security**: Security-related decisions and standards

## Usage

### Viewing ADRs

- **Summary**: See [SUMMARY.md](data/SUMMARY.md) for an overview of all ADRs
- **Individual ADRs**: Browse individual ADR files in [data/](data/)
- **Timeline**: View decisions chronologically in the summary

### Regenerating ADRs

To regenerate ADRs from git history:

```bash
# From project root
just generate-adrs

# Or manually
cd scripts && go run generate-adrs.go
```

### Focused Analysis

For focused recommendations of significant architectural decisions:

```bash
# From project root
just focused-adr-analysis

# Or manually
cd scripts && go build -o focused-adr-analysis focused-adr-analysis.go && ./focused-adr-analysis
```

### Comprehensive Analysis

For complete analysis of all potential architectural decisions:

```bash
# From project root
just analyze-git-comprehensive

# Or manually
cd scripts && go build -o analyze-git-adrs analyze-git-adrs.go && ./analyze-git-adrs
```

### Adding New ADRs

To add new ADRs:

1. Update the `scripts/generate-adrs.go` file
2. Add new ADR definitions to the appropriate category function
3. Run the generation script
4. Review and refine the generated content

## ADR Format

Each ADR follows this structure:

```markdown
# ADR-XXX: Title

**Date:** YYYY-MM-DD
**Author:** Author Name
**Status:** Accepted/Deprecated/Superseded

## Context

Problem description and background.

## Decision

What was decided and why.

## Consequences

- ✅ Positive consequences
- ❌ Negative consequences

## Evidence

- Supporting evidence from git history
- Code patterns and implementations

## Related Files

- `path/to/relevant/files`
```

## Analysis

These ADRs provide insights into:

- **Architectural Evolution**: How the system architecture has evolved over time
- **Decision Patterns**: Common patterns in architectural decision-making
- **Trade-offs**: Understanding of the trade-offs made in different decisions
- **Technical Debt**: Areas where decisions may have introduced complexity

## Integration with Development

ADRs help with:

- **Onboarding**: New developers can understand architectural decisions
- **Refactoring**: Understanding why certain patterns were chosen
- **Planning**: Making informed decisions about future changes
- **Documentation**: Maintaining architectural knowledge
