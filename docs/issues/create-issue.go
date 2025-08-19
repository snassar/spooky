package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type IssueRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type IssueResponse struct {
	Number int `json:"number"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run create-issue.go <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found\n", filename)
		os.Exit(1)
	}

	// Check environment variables
	token := os.Getenv("CODEBERG_TOKEN")
	owner := os.Getenv("CODEBERG_OWNER")
	repo := os.Getenv("CODEBERG_REPO")

	if token == "" {
		fmt.Println("Error: CODEBERG_TOKEN environment variable not set")
		os.Exit(1)
	}
	if owner == "" {
		fmt.Println("Error: CODEBERG_OWNER environment variable not set")
		os.Exit(1)
	}
	if repo == "" {
		fmt.Println("Error: CODEBERG_REPO environment variable not set")
		os.Exit(1)
	}

	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Split(string(content), "\n")

	// Check if issue number already exists
	issueNumberRegex := regexp.MustCompile(`^#\s*(\d+)`)
	for _, line := range lines {
		if match := issueNumberRegex.FindStringSubmatch(line); match != nil {
			fmt.Printf("Issue number already exists in file: #%s\n", match[1])
			return
		}
	}

	// Extract title (first line starting with #)
	var title string
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			break
		}
	}

	if title == "" {
		fmt.Println("Error: No title found (first line should start with '# ')")
		os.Exit(1)
	}

	// Extract body (everything after first line, with # converted to ##)
	var bodyLines []string
	foundFirst := false
	for _, line := range lines {
		if !foundFirst && strings.HasPrefix(line, "# ") {
			foundFirst = true
			continue
		}
		if foundFirst {
			if strings.HasPrefix(line, "# ") {
				bodyLines = append(bodyLines, "## "+strings.TrimPrefix(line, "# "))
			} else {
				bodyLines = append(bodyLines, line)
			}
		}
	}

	body := strings.Join(bodyLines, "\n")

	fmt.Printf("Creating new issue: %s\n", title)

	// Create issue request
	issueReq := IssueRequest{
		Title: title,
		Body:  body,
	}

	jsonData, err := json.Marshal(issueReq)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}

	// Make HTTP request
	url := fmt.Sprintf("https://codeberg.org/api/v1/repos/%s/%s/issues", owner, repo)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	// Parse response
	var issueResp IssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Issue created successfully on Codeberg: #%d\n", issueResp.Number)

	// Add issue number to file
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file for update: %v\n", err)
		os.Exit(1)
	}

	newContent := fmt.Sprintf("# %d\n%s", issueResp.Number, string(fileContent))
	if err := os.WriteFile(filename, []byte(newContent), 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s with issue number #%d\n", filename, issueResp.Number)
}
