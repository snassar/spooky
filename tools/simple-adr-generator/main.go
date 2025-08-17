package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type GitCommit struct {
	Hash      string
	Author    string
	Date      time.Time
	Message   string
	Files     []string
	Additions int
	Deletions int
}

type ADR struct {
	ID           string
	Title        string
	Date         time.Time
	Author       string
	Context      string
	Decision     string
	Consequences []string
	Evidence     []string
	Status       string
	Files        []string
	CommitHash   string
}

func main() {
	// Create output directory
	outputDir := "../docs/adr"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Get git history
	commits := getGitHistory()

	// Generate ADRs from significant commits
	adrs := generateADRs(commits)

	// Write ADR files
	for _, adr := range adrs {
		writeADRFile(outputDir, adr)
	}

	// Write summary
	writeSummary(outputDir, adrs)

	fmt.Printf("Generated %d ADRs in %s/\n", len(adrs), outputDir)
}

func getGitHistory() []GitCommit {
	output, err := exec.Command("git", "log", "--stat", "--oneline", "--format=%H|%an|%ai|%s").Output()
	if err != nil {
		fmt.Printf("Error getting git history: %v\n", err)
		return nil
	}

	var commits []GitCommit
	lines := strings.Split(string(output), "\n")

	var currentCommit GitCommit
	for _, line := range lines {
		if strings.HasPrefix(line, " ") && strings.Contains(line, "|") {
			// Parse file stats
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				stats := strings.TrimSpace(parts[1])
				if strings.Contains(stats, "insertions") {
					re := regexp.MustCompile(`(\d+) insertions?.*?(\d+) deletions?`)
					matches := re.FindStringSubmatch(stats)
					if len(matches) >= 3 {
						var additions, deletions int
						fmt.Sscanf(matches[1], "%d", &additions)
						fmt.Sscanf(matches[2], "%d", &deletions)

						currentCommit.Files = append(currentCommit.Files, parts[0])
						currentCommit.Additions += additions
						currentCommit.Deletions += deletions
					}
				}
			}
		} else if len(line) > 0 && !strings.HasPrefix(line, " ") {
			// New commit
			if currentCommit.Hash != "" {
				commits = append(commits, currentCommit)
			}

			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[2])
				currentCommit = GitCommit{
					Hash:    parts[0],
					Author:  parts[1],
					Date:    date,
					Message: parts[3],
				}
			}
		}
	}

	if currentCommit.Hash != "" {
		commits = append(commits, currentCommit)
	}

	return commits
}

func generateADRs(commits []GitCommit) []ADR {
	var adrs []ADR
	adrCounter := 1

	// Keywords that indicate architectural decisions
	keywords := map[string]float64{
		"architect": 0.8, "architecture": 0.8, "design": 0.6,
		"refactor": 0.5, "restructure": 0.7, "redesign": 0.7,
		"interface": 0.6, "system": 0.5, "breaking": 0.9,
		"major": 0.7, "consolidate": 0.6, "unify": 0.6,
		"standardize": 0.5, "migration": 0.8, "deprecate": 0.7,
	}

	for _, commit := range commits {
		confidence := 0.0
		messageLower := strings.ToLower(commit.Message)

		// Calculate confidence based on keywords
		for keyword, score := range keywords {
			if strings.Contains(messageLower, keyword) {
				confidence += score
			}
		}

		// Large changes
		if commit.Additions > 1000 || commit.Deletions > 500 {
			confidence += 0.3
		}

		// Architectural files
		architecturalFiles := countArchitecturalFiles(commit.Files)
		confidence += float64(architecturalFiles) * 0.1

		// Only generate ADRs for significant commits
		if confidence >= 0.5 {
			adr := ADR{
				ID:           fmt.Sprintf("ADR-%03d", adrCounter),
				Title:        commit.Message,
				Date:         commit.Date,
				Author:       commit.Author,
				Context:      generateContext(commit),
				Decision:     generateDecision(commit),
				Consequences: generateConsequences(commit),
				Evidence:     generateEvidence(commit),
				Status:       "Accepted",
				Files:        commit.Files,
				CommitHash:   commit.Hash,
			}
			adrs = append(adrs, adr)
			adrCounter++
		}
	}

	// Sort by date (newest first)
	sort.Slice(adrs, func(i, j int) bool {
		return adrs[i].Date.After(adrs[j].Date)
	})

	return adrs
}

func countArchitecturalFiles(files []string) int {
	count := 0
	architecturalPatterns := []string{
		"internal/interfaces/", "internal/types/", "internal/schemas/",
		"internal/actions/", "internal/facts/", "internal/machines/",
		"internal/ssh/", "internal/templates/", "internal/variables/",
		"docs/design/", "docs/plans/", ".cursor/rules/",
	}

	for _, file := range files {
		for _, pattern := range architecturalPatterns {
			if strings.Contains(file, pattern) {
				count++
				break
			}
		}
	}

	return count
}

func generateContext(commit GitCommit) string {
	return fmt.Sprintf("This decision was made during development of the spooky project. The commit %s by %s on %s represents a significant architectural change.",
		commit.Hash[:8], commit.Author, commit.Date.Format("2006-01-02"))
}

func generateDecision(commit GitCommit) string {
	return fmt.Sprintf("Based on the commit message '%s', this represents a decision to %s.",
		commit.Message, extractDecisionFromMessage(commit.Message))
}

func extractDecisionFromMessage(message string) string {
	messageLower := strings.ToLower(message)

	if strings.Contains(messageLower, "architect") {
		return "implement architectural changes"
	}
	if strings.Contains(messageLower, "refactor") {
		return "refactor existing code structure"
	}
	if strings.Contains(messageLower, "interface") {
		return "modify interface definitions"
	}
	if strings.Contains(messageLower, "system") {
		return "change system behavior"
	}
	if strings.Contains(messageLower, "breaking") {
		return "introduce breaking changes"
	}

	return "make significant changes to the codebase"
}

func generateConsequences(commit GitCommit) []string {
	var consequences []string

	if commit.Additions > 1000 || commit.Deletions > 500 {
		consequences = append(consequences, "Large codebase changes requiring careful review")
	}

	architecturalFiles := countArchitecturalFiles(commit.Files)
	if architecturalFiles > 0 {
		consequences = append(consequences, fmt.Sprintf("Changes to %d architectural files", architecturalFiles))
	}

	consequences = append(consequences, "May require updates to dependent code")
	consequences = append(consequences, "Should be reviewed for backward compatibility")

	return consequences
}

func generateEvidence(commit GitCommit) []string {
	var evidence []string

	evidence = append(evidence, fmt.Sprintf("Commit: %s", commit.Hash))
	evidence = append(evidence, fmt.Sprintf("Author: %s", commit.Author))
	evidence = append(evidence, fmt.Sprintf("Date: %s", commit.Date.Format("2006-01-02 15:04:05")))
	evidence = append(evidence, fmt.Sprintf("Changes: +%d -%d lines", commit.Additions, commit.Deletions))

	if len(commit.Files) > 0 {
		evidence = append(evidence, fmt.Sprintf("Files affected: %s", strings.Join(commit.Files, ", ")))
	}

	return evidence
}

func writeADRFile(outputDir string, adr ADR) {
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.md", adr.ID))
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# %s: %s\n\n", adr.ID, adr.Title)
	fmt.Fprintf(writer, "**Date:** %s  \n", adr.Date.Format("2006-01-02"))
	fmt.Fprintf(writer, "**Author:** %s  \n", adr.Author)
	fmt.Fprintf(writer, "**Status:** %s  \n\n", adr.Status)

	fmt.Fprintf(writer, "## Context\n\n%s\n\n", adr.Context)

	fmt.Fprintf(writer, "## Decision\n\n%s\n\n", adr.Decision)

	if len(adr.Consequences) > 0 {
		fmt.Fprintf(writer, "## Consequences\n\n")
		for _, consequence := range adr.Consequences {
			fmt.Fprintf(writer, "- %s\n", consequence)
		}
		fmt.Fprintf(writer, "\n")
	}

	if len(adr.Evidence) > 0 {
		fmt.Fprintf(writer, "## Evidence\n\n")
		for _, evidence := range adr.Evidence {
			fmt.Fprintf(writer, "- %s\n", evidence)
		}
		fmt.Fprintf(writer, "\n")
	}
}

func writeSummary(outputDir string, adrs []ADR) {
	filename := filepath.Join(outputDir, "SUMMARY.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating summary file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Architecture Decision Records Summary\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "Total ADRs: %d\n\n", len(adrs))

	fmt.Fprintf(writer, "## ADRs by Date\n\n")
	for _, adr := range adrs {
		fmt.Fprintf(writer, "- [%s](%s.md) - %s (%s)\n",
			adr.ID, adr.ID, adr.Title, adr.Date.Format("2006-01-02"))
	}
}
