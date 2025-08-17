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

type PotentialADR struct {
	Commit     GitCommit
	Category   string
	Confidence float64
	Reason     string
}

func main() {
	// Create output directory
	outputDir := "../docs/adr/analysis"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Analyze git history
	commits := getGitHistory()
	potentialADRs := analyzeCommits(commits)

	// Write analysis report
	writeAnalysisReport(outputDir, commits, potentialADRs)

	fmt.Printf("Analysis complete. Found %d commits, %d potential ADRs\n", len(commits), len(potentialADRs))
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
						if _, err := fmt.Sscanf(matches[1], "%d", &additions); err != nil {
							// Log error but continue
							fmt.Printf("Warning: failed to parse additions: %v\n", err)
						}
						if _, err := fmt.Sscanf(matches[2], "%d", &deletions); err != nil {
							// Log error but continue
							fmt.Printf("Warning: failed to parse deletions: %v\n", err)
						}

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

func analyzeCommits(commits []GitCommit) []PotentialADR {
	var potentialADRs []PotentialADR

	for _, commit := range commits {
		confidence := calculateConfidence(commit)
		if confidence > 0.3 {
			category := determineCategory(commit)
			reason := determineReason(commit)

			potentialADRs = append(potentialADRs, PotentialADR{
				Commit:     commit,
				Category:   category,
				Confidence: confidence,
				Reason:     reason,
			})
		}
	}

	// Sort by confidence
	sort.Slice(potentialADRs, func(i, j int) bool {
		return potentialADRs[i].Confidence > potentialADRs[j].Confidence
	})

	return potentialADRs
}

func calculateConfidence(commit GitCommit) float64 {
	confidence := 0.0

	// Keywords that indicate architectural decisions
	keywords := map[string]float64{
		"architect": 0.8, "architecture": 0.8, "design": 0.6,
		"refactor": 0.5, "restructure": 0.7, "redesign": 0.7,
		"interface": 0.6, "system": 0.5, "breaking": 0.9,
		"major": 0.7, "consolidate": 0.6, "unify": 0.6,
		"standardize": 0.5, "migration": 0.8, "deprecate": 0.7,
	}

	messageLower := strings.ToLower(commit.Message)
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

	return confidence
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

func determineCategory(commit GitCommit) string {
	messageLower := strings.ToLower(commit.Message)

	if strings.Contains(messageLower, "architect") || strings.Contains(messageLower, "design") {
		return "Architecture"
	}
	if strings.Contains(messageLower, "interface") {
		return "Interface"
	}
	if strings.Contains(messageLower, "system") {
		return "System"
	}
	if strings.Contains(messageLower, "breaking") || strings.Contains(messageLower, "migration") {
		return "Breaking Change"
	}
	if strings.Contains(messageLower, "refactor") {
		return "Refactoring"
	}
	if strings.Contains(messageLower, "feat") {
		return "Features"
	}
	if strings.Contains(messageLower, "performance") || strings.Contains(messageLower, "optimize") {
		return "Performance"
	}
	if strings.Contains(messageLower, "validation") || strings.Contains(messageLower, "schema") {
		return "Configuration"
	}

	return "Other"
}

func determineReason(commit GitCommit) string {
	var reasons []string

	messageLower := strings.ToLower(commit.Message)

	// Check for keywords
	keywords := []string{"architect", "design", "refactor", "interface", "system", "breaking", "major", "consolidate", "unify", "standardize", "migration", "deprecate"}
	for _, keyword := range keywords {
		if strings.Contains(messageLower, keyword) {
			reasons = append(reasons, fmt.Sprintf("Contains keyword: %s", keyword))
		}
	}

	// Check for large changes
	if commit.Additions > 1000 || commit.Deletions > 500 {
		reasons = append(reasons, fmt.Sprintf("Large change: %d lines", commit.Additions+commit.Deletions))
	}

	// Check for architectural files
	architecturalFiles := countArchitecturalFiles(commit.Files)
	if architecturalFiles > 0 {
		reasons = append(reasons, fmt.Sprintf("Architectural files: %d", architecturalFiles))
	}

	return strings.Join(reasons, "; ")
}

func writeAnalysisReport(outputDir string, commits []GitCommit, potentialADRs []PotentialADR) {
	filename := filepath.Join(outputDir, "git-analysis.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Git History Analysis for ADR Discovery\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "## Summary\n\n")
	fmt.Fprintf(writer, "- **Total commits analyzed:** %d\n", len(commits))
	fmt.Fprintf(writer, "- **Potential ADRs found:** %d\n", len(potentialADRs))
	if len(commits) > 0 {
		fmt.Fprintf(writer, "- **Analysis period:** %s to %s\n\n",
			commits[len(commits)-1].Date.Format("2006-01-02"),
			commits[0].Date.Format("2006-01-02"))
	}

	// Group by confidence
	highConfidence := filterByConfidence(potentialADRs, 0.7)
	mediumConfidence := filterByConfidence(potentialADRs, 0.5)
	mediumConfidence = filterByConfidence(mediumConfidence, 0.69) // Exclude high confidence

	if len(highConfidence) > 0 {
		fmt.Fprintf(writer, "## High Confidence ADRs (Confidence >= 0.7)\n\n")
		for _, adr := range highConfidence {
			writeADREntry(writer, adr)
		}
	}

	if len(mediumConfidence) > 0 {
		fmt.Fprintf(writer, "## Medium Confidence ADRs (Confidence 0.5-0.69)\n\n")
		for _, adr := range mediumConfidence {
			writeADREntry(writer, adr)
		}
	}

	// Group by category
	fmt.Fprintf(writer, "## By Category\n\n")
	categories := groupByCategory(potentialADRs)
	for category, adrs := range categories {
		fmt.Fprintf(writer, "### %s (%d ADRs)\n\n", category, len(adrs))
		for _, adr := range adrs {
			fmt.Fprintf(writer, "- [%s](%s) (%.2f) - %s\n",
				adr.Commit.Message, adr.Commit.Hash[:8], adr.Confidence, adr.Commit.Date.Format("2006-01-02"))
		}
		fmt.Fprintf(writer, "\n")
	}

	fmt.Fprintf(writer, "## Recommendations\n\n")
	fmt.Fprintf(writer, "### High Priority ADRs to Document\n\n")
	highPriorityCount := 0
	for _, adr := range potentialADRs {
		if adr.Confidence >= 0.7 && highPriorityCount < 10 {
			fmt.Fprintf(writer, "1. **%s** - %s (%.2f confidence)\n",
				adr.Commit.Hash[:8], adr.Commit.Message, adr.Confidence)
			highPriorityCount++
		}
	}
}

func writeADREntry(writer *bufio.Writer, adr PotentialADR) {
	fmt.Fprintf(writer, "### %s\n\n", adr.Commit.Message)
	fmt.Fprintf(writer, "- **Commit:** %s\n", adr.Commit.Hash[:8])
	fmt.Fprintf(writer, "- **Date:** %s\n", adr.Commit.Date.Format("2006-01-02"))
	fmt.Fprintf(writer, "- **Author:** %s\n", adr.Commit.Author)
	fmt.Fprintf(writer, "- **Category:** %s\n", adr.Category)
	fmt.Fprintf(writer, "- **Confidence:** %.2f\n", adr.Confidence)
	fmt.Fprintf(writer, "- **Reason:** %s\n", adr.Reason)
	fmt.Fprintf(writer, "- **Changes:** +%d -%d\n", adr.Commit.Additions, adr.Commit.Deletions)
	fmt.Fprintf(writer, "- **Files:** %s\n\n", strings.Join(adr.Commit.Files, ", "))
}

func filterByConfidence(adrs []PotentialADR, minConfidence float64) []PotentialADR {
	var filtered []PotentialADR
	for _, adr := range adrs {
		if adr.Confidence >= minConfidence {
			filtered = append(filtered, adr)
		}
	}
	return filtered
}

func groupByCategory(adrs []PotentialADR) map[string][]PotentialADR {
	categories := make(map[string][]PotentialADR)
	for _, adr := range adrs {
		categories[adr.Category] = append(categories[adr.Category], adr)
	}
	return categories
}
