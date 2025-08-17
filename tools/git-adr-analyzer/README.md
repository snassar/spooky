# Git ADR Analyzer Tool

## Purpose
Analyzes ADRs (Architecture Decision Records) from Git history to track decision evolution.

## Usage
```bash
# Build the tool
just build

# Analyze ADRs from Git history
just run
```

## Output
Analyzes Git history to understand how ADRs have evolved over time and identifies patterns in decision-making.

## Configuration
This tool scans Git history for ADR-related commits and analyzes decision patterns.

## Dependencies
- Go 1.24+
- Git repository access
- Access to the spooky codebase
