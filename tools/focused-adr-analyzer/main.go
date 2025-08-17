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

type SignificantCommit struct {
	Hash      string
	Author    string
	Date      time.Time
	Message   string
	Category  string
	Impact    string
	Reason    string
	Files     []string
	Additions int
	Deletions int
}

type ADRRecommendation struct {
	Commit         SignificantCommit
	Priority       string // "High", "Medium", "Low"
	Reason         string
	SuggestedTitle string
}

func main() {
	// Create output directory
	outputDir := "../docs/adr/recommendations"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Get significant commits
	significantCommits := getSignificantCommits()
	recommendations := analyzeForADRs(significantCommits)

	// Write focused recommendations
	writeFocusedRecommendations(outputDir, recommendations)

	fmt.Printf("Analysis complete. Found %d significant commits, %d ADR recommendations\n",
		len(significantCommits), len(recommendations))
}

func getSignificantCommits() []SignificantCommit {
	// Look for commits with specific patterns that indicate architectural decisions
	patterns := []string{
		"refactor.*architecture",
		"refactor.*interface",
		"refactor.*system",
		"feat.*architecture",
		"feat.*system",
		"breaking.*change",
		"major.*refactor",
		"restructure",
		"redesign",
		"migration",
		"deprecate",
		"replace.*with",
		"consolidate",
		"unify",
		"standardize",
	}

	var commits []SignificantCommit

	for _, pattern := range patterns {
		output, err := exec.Command("git", "log", "--grep", pattern, "--oneline", "--stat").Output()
		if err != nil {
			continue
		}

		commits = append(commits, parseSignificantCommits(string(output), pattern)...)
	}

	// Also look for large commits that might be architectural
	output, err := exec.Command("git", "log", "--stat", "--oneline").Output()
	if err == nil {
		commits = append(commits, findLargeArchitecturalCommits(string(output))...)
	}

	// Remove duplicates and sort by date
	uniqueCommits := deduplicateCommits(commits)
	sort.Slice(uniqueCommits, func(i, j int) bool {
		return uniqueCommits[i].Date.After(uniqueCommits[j].Date)
	})

	return uniqueCommits
}

func parseSignificantCommits(output, pattern string) []SignificantCommit {
	var commits []SignificantCommit
	lines := strings.Split(output, "\n")

	var currentCommit SignificantCommit
	for _, line := range lines {
		if strings.HasPrefix(line, " ") && strings.Contains(line, "|") {
			// Parse file stats
			parts := strings.Split(line, "|")
			if len(parts) >= 2 {
				stats := strings.TrimSpace(parts[1])
				if strings.Contains(stats, "insertions") {
					// Extract numbers
					re := regexp.MustCompile(`(\d+) insertions?.*?(\d+) deletions?`)
					matches := re.FindStringSubmatch(stats)
					if len(matches) >= 3 {
						fmt.Sscanf(matches[1], "%d", &currentCommit.Additions)
						fmt.Sscanf(matches[2], "%d", &currentCommit.Deletions)
					}
				}
			}
		} else if len(line) > 0 && !strings.HasPrefix(line, " ") {
			// New commit
			if currentCommit.Hash != "" {
				// Get commit date
				if date, err := getCommitDate(currentCommit.Hash); err == nil {
					currentCommit.Date = date
				}
				commits = append(commits, currentCommit)
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) >= 2 {
				currentCommit = SignificantCommit{
					Hash:     parts[0],
					Message:  parts[1],
					Category: getCategoryFromPattern(pattern),
					Reason:   fmt.Sprintf("Matches pattern: %s", pattern),
					Date:     time.Now(), // Default to now, will be updated
				}
			}
		}
	}

	if currentCommit.Hash != "" {
		// Get commit date for last commit
		if date, err := getCommitDate(currentCommit.Hash); err == nil {
			currentCommit.Date = date
		}
		commits = append(commits, currentCommit)
	}

	return commits
}

func findLargeArchitecturalCommits(output string) []SignificantCommit {
	var commits []SignificantCommit
	lines := strings.Split(output, "\n")

	var currentCommit SignificantCommit
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

						// Only consider very large commits
						if additions > 1000 || deletions > 500 {
							currentCommit.Additions = additions
							currentCommit.Deletions = deletions
							currentCommit.Category = "Large Change"
							currentCommit.Reason = fmt.Sprintf("Large commit: +%d -%d lines", additions, deletions)
						}
					}
				}
			}
		} else if len(line) > 0 && !strings.HasPrefix(line, " ") {
			// New commit
			if currentCommit.Hash != "" && (currentCommit.Additions > 1000 || currentCommit.Deletions > 500) {
				commits = append(commits, currentCommit)
			}

			parts := strings.SplitN(line, " ", 2)
			if len(parts) >= 2 {
				currentCommit = SignificantCommit{
					Hash:    parts[0],
					Message: parts[1],
				}
			}
		}
	}

	if currentCommit.Hash != "" && (currentCommit.Additions > 1000 || currentCommit.Deletions > 500) {
		commits = append(commits, currentCommit)
	}

	return commits
}

func getCategoryFromPattern(pattern string) string {
	switch {
	case strings.Contains(pattern, "architecture"):
		return "Architecture"
	case strings.Contains(pattern, "interface"):
		return "Interface"
	case strings.Contains(pattern, "system"):
		return "System"
	case strings.Contains(pattern, "breaking"):
		return "Breaking Change"
	case strings.Contains(pattern, "migration"):
		return "Migration"
	default:
		return "Refactoring"
	}
}

func getCommitDate(hash string) (time.Time, error) {
	output, err := exec.Command("git", "log", "-1", "--format=%ci", hash).Output()
	if err != nil {
		return time.Time{}, err
	}

	dateStr := strings.TrimSpace(string(output))
	return time.Parse("2006-01-02 15:04:05 -0700", dateStr)
}

func deduplicateCommits(commits []SignificantCommit) []SignificantCommit {
	seen := make(map[string]bool)
	var unique []SignificantCommit

	for _, commit := range commits {
		if !seen[commit.Hash] {
			seen[commit.Hash] = true
			unique = append(unique, commit)
		}
	}

	return unique
}

func analyzeForADRs(commits []SignificantCommit) []ADRRecommendation {
	var recommendations []ADRRecommendation

	for _, commit := range commits {
		priority := determinePriority(commit)
		if priority != "Low" {
			recommendations = append(recommendations, ADRRecommendation{
				Commit:         commit,
				Priority:       priority,
				Reason:         commit.Reason,
				SuggestedTitle: generateTitle(commit),
			})
		}
	}

	return recommendations
}

func determinePriority(commit SignificantCommit) string {
	// High priority: Breaking changes, major refactors, architecture changes
	if strings.Contains(strings.ToLower(commit.Message), "breaking") ||
		strings.Contains(strings.ToLower(commit.Message), "major") ||
		strings.Contains(strings.ToLower(commit.Message), "architecture") ||
		commit.Additions > 5000 || commit.Deletions > 2000 {
		return "High"
	}

	// Medium priority: System changes, interface changes, large refactors
	if strings.Contains(strings.ToLower(commit.Message), "system") ||
		strings.Contains(strings.ToLower(commit.Message), "interface") ||
		strings.Contains(strings.ToLower(commit.Message), "refactor") ||
		commit.Additions > 1000 || commit.Deletions > 500 {
		return "Medium"
	}

	return "Low"
}

func generateTitle(commit SignificantCommit) string {
	// Extract meaningful title from commit message
	message := commit.Message

	// Remove common prefixes
	prefixes := []string{"feat:", "fix:", "refactor:", "docs:", "test:", "chore:"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(message, prefix) {
			message = strings.TrimSpace(strings.TrimPrefix(message, prefix))
			break
		}
	}

	// Capitalize first letter
	if len(message) > 0 {
		message = strings.ToUpper(string(message[0])) + message[1:]
	}

	return message
}

func writeFocusedRecommendations(outputDir string, recommendations []ADRRecommendation) {
	filename := filepath.Join(outputDir, "focused-adr-recommendations.md")
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "# Focused ADR Recommendations\n\n")
	fmt.Fprintf(writer, "Generated on: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(writer, "## Summary\n\n")
	fmt.Fprintf(writer, "- **Total recommendations:** %d\n", len(recommendations))
	fmt.Fprintf(writer, "- **High priority:** %d\n", countByPriority(recommendations, "High"))
	fmt.Fprintf(writer, "- **Medium priority:** %d\n", countByPriority(recommendations, "Medium"))
	fmt.Fprintf(writer, "- **Low priority:** %d\n\n", countByPriority(recommendations, "Low"))

	// Group by priority
	highPriority := filterByPriority(recommendations, "High")
	mediumPriority := filterByPriority(recommendations, "Medium")

	if len(highPriority) > 0 {
		fmt.Fprintf(writer, "## High Priority ADRs\n\n")
		fmt.Fprintf(writer, "These commits represent significant architectural decisions that should be documented:\n\n")

		for _, rec := range highPriority {
			fmt.Fprintf(writer, "### %s\n\n", rec.SuggestedTitle)
			hash := rec.Commit.Hash
			if len(hash) >= 8 {
				hash = hash[:8]
			}
			fmt.Fprintf(writer, "- **Commit:** %s\n", hash)
			fmt.Fprintf(writer, "- **Date:** %s\n", rec.Commit.Date.Format("2006-01-02"))
			fmt.Fprintf(writer, "- **Category:** %s\n", rec.Commit.Category)
			fmt.Fprintf(writer, "- **Impact:** +%d -%d lines\n", rec.Commit.Additions, rec.Commit.Deletions)
			fmt.Fprintf(writer, "- **Reason:** %s\n\n", rec.Reason)
		}
	}

	if len(mediumPriority) > 0 {
		fmt.Fprintf(writer, "## Medium Priority ADRs\n\n")
		fmt.Fprintf(writer, "These commits may represent architectural decisions worth documenting:\n\n")

		for _, rec := range mediumPriority {
			fmt.Fprintf(writer, "### %s\n\n", rec.SuggestedTitle)
			hash := rec.Commit.Hash
			if len(hash) >= 8 {
				hash = hash[:8]
			}
			fmt.Fprintf(writer, "- **Commit:** %s\n", hash)
			fmt.Fprintf(writer, "- **Date:** %s\n", rec.Commit.Date.Format("2006-01-02"))
			fmt.Fprintf(writer, "- **Category:** %s\n", rec.Commit.Category)
			fmt.Fprintf(writer, "- **Impact:** +%d -%d lines\n", rec.Commit.Additions, rec.Commit.Deletions)
			fmt.Fprintf(writer, "- **Reason:** %s\n\n", rec.Reason)
		}
	}

	fmt.Fprintf(writer, "## Next Steps\n\n")
	fmt.Fprintf(writer, "1. Review high priority recommendations\n")
	fmt.Fprintf(writer, "2. Create ADRs for the most significant decisions\n")
	fmt.Fprintf(writer, "3. Update existing ADRs if decisions have evolved\n")
	fmt.Fprintf(writer, "4. Consider documenting medium priority items as needed\n\n")

	fmt.Fprintf(writer, "## Methodology\n\n")
	fmt.Fprintf(writer, "This analysis focuses on:\n")
	fmt.Fprintf(writer, "- Commits with architectural keywords (refactor, architecture, system, etc.)\n")
	fmt.Fprintf(writer, "- Large commits (>1000 additions or >500 deletions)\n")
	fmt.Fprintf(writer, "- Breaking changes and major refactors\n")
	fmt.Fprintf(writer, "- Interface and system-level changes\n")
}

func countByPriority(recommendations []ADRRecommendation, priority string) int {
	count := 0
	for _, rec := range recommendations {
		if rec.Priority == priority {
			count++
		}
	}
	return count
}

func filterByPriority(recommendations []ADRRecommendation, priority string) []ADRRecommendation {
	var filtered []ADRRecommendation
	for _, rec := range recommendations {
		if rec.Priority == priority {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}
