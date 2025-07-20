package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"spooky/internal/facts"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// TestCLICommand is a helper struct for testing CLI commands
type TestCLICommand struct {
	cmd    *cobra.Command
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

// NewTestCLICommand creates a new test CLI command
func NewTestCLICommand(cmd *cobra.Command) *TestCLICommand {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return &TestCLICommand{
		cmd:    cmd,
		stdout: stdout,
		stderr: stderr,
	}
}

// Execute runs the command with the given arguments
func (t *TestCLICommand) Execute(args ...string) error {
	t.cmd.SetArgs(args)
	return t.cmd.Execute()
}

// GetOutput returns the stdout content
func (t *TestCLICommand) GetOutput() string {
	return t.stdout.String()
}

// GetError returns the stderr content
func (t *TestCLICommand) GetError() string {
	return t.stderr.String()
}

// setupTestEnvironment creates a temporary test environment
func setupTestEnvironment(t *testing.T) (string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "spooky-facts-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Set environment variables
	originalPath := os.Getenv("SPOOKY_FACTS_PATH")
	originalFormat := os.Getenv("SPOOKY_FACTS_FORMAT")

	os.Setenv("SPOOKY_FACTS_PATH", filepath.Join(tempDir, "test.db"))
	os.Setenv("SPOOKY_FACTS_FORMAT", "badgerdb")

	// Cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
		os.Setenv("SPOOKY_FACTS_PATH", originalPath)
		os.Setenv("SPOOKY_FACTS_FORMAT", originalFormat)
	}

	return tempDir, cleanup
}

func TestFactsCollectCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{
		Use:   "facts collect",
		Short: "Collect facts from a server",
		RunE:  runFactsCollect,
	}

	testCmd := NewTestCLICommand(cmd)

	// Test collecting all facts from local machine
	err := testCmd.Execute("local")
	if err != nil {
		t.Fatalf("Failed to execute facts collect command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts collect command executed successfully")
}

func TestFactsGetCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{
		Use:   "facts get",
		Short: "Get a specific fact",
		RunE:  runFactsGet,
	}

	testCmd := NewTestCLICommand(cmd)

	// Test getting a specific fact
	err := testCmd.Execute("local", "hostname")
	if err != nil {
		t.Fatalf("Failed to execute facts get command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts get command executed successfully")
}

func TestFactsListCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First, collect some facts to have data to list
	collectCmd := &cobra.Command{
		Use:   "facts collect",
		Short: "Collect facts from a server",
		RunE:  runFactsCollect,
	}

	collectTestCmd := NewTestCLICommand(collectCmd)
	err := collectTestCmd.Execute("local")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Now test listing facts
	listCmd := &cobra.Command{
		Use:   "facts list",
		Short: "List stored facts",
		RunE:  runFactsList,
	}

	listTestCmd := NewTestCLICommand(listCmd)
	err = listTestCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute facts list command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts list command executed successfully")
}

func TestFactsGatherCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test command
	cmd := &cobra.Command{
		Use:   "facts gather",
		Short: "Gather facts from multiple hosts",
		RunE:  runFactsGather,
	}

	testCmd := NewTestCLICommand(cmd)

	// Test gathering facts from local machine
	err := testCmd.Execute("local")
	if err != nil {
		t.Fatalf("Failed to execute facts gather command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts gather command executed successfully")
}

func TestFactsQueryCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First, collect some facts to have data to query
	collectCmd := &cobra.Command{
		Use:   "facts collect",
		Short: "Collect facts from a server",
		RunE:  runFactsCollect,
	}

	collectTestCmd := NewTestCLICommand(collectCmd)
	err := collectTestCmd.Execute("local")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Now test querying facts
	queryCmd := &cobra.Command{
		Use:   "facts query",
		Short: "Query stored facts",
		RunE:  runFactsQuery,
	}

	queryTestCmd := NewTestCLICommand(queryCmd)

	// Test querying with machine name
	err = queryTestCmd.Execute("machine=local")
	if err != nil {
		t.Fatalf("Failed to execute facts query command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts query command executed successfully")
}

func TestFactsExportCommand(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// First, collect some facts to have data to export
	collectCmd := &cobra.Command{
		Use:   "facts collect",
		Short: "Collect facts from a server",
		RunE:  runFactsCollect,
	}

	collectTestCmd := NewTestCLICommand(collectCmd)
	err := collectTestCmd.Execute("local")
	if err != nil {
		t.Fatalf("Failed to collect facts: %v", err)
	}

	// Create export file path
	exportFile := filepath.Join(tempDir, "export.json")

	// Test exporting facts
	exportCmd := &cobra.Command{
		Use:   "facts export",
		Short: "Export facts to JSON",
		RunE:  runFactsExport,
	}

	// Add the output flag
	exportCmd.Flags().StringVar(&factsOutput, "output", "", "Output file path (default: stdout)")

	exportTestCmd := NewTestCLICommand(exportCmd)
	err = exportTestCmd.Execute("--output", exportFile)
	if err != nil {
		t.Fatalf("Failed to execute facts export command: %v", err)
	}

	// Verify export file was created
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		t.Error("Export file was not created")
	}

	// Verify export file contains valid JSON
	exportData, err := os.ReadFile(exportFile)
	if err != nil {
		t.Fatalf("Failed to read export file: %v", err)
	}

	var exportedFacts map[string]*facts.MachineFacts
	err = json.Unmarshal(exportData, &exportedFacts)
	if err != nil {
		t.Fatalf("Failed to parse exported JSON: %v", err)
	}

	if len(exportedFacts) == 0 {
		t.Error("Expected exported facts to contain data")
	}
}

func TestFactsImportCommand(t *testing.T) {
	tempDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test data to import (map format expected by ImportFromJSON)
	testFacts := map[string]*facts.MachineFacts{
		"test-machine-1": {
			MachineID:   "test-machine-1",
			MachineName: "test-machine-1",
			Hostname:    "test-host-1",
			OS:          "Linux",
			OSVersion:   "Ubuntu 22.04",
			IPAddresses: []string{"192.168.1.100"},
			PrimaryIP:   "192.168.1.100",
			CPU: facts.CPUInfo{
				Cores:     4,
				Model:     "Intel i7",
				Arch:      "x86_64",
				Frequency: "2.4 GHz",
			},
			Memory: facts.MemoryInfo{
				Total:     8589934592,
				Used:      4294967296,
				Available: 4294967296,
			},
			SystemID:  "test-system-id-1",
			Tags:      map[string]string{"environment": "test"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Create import file
	importFile := filepath.Join(tempDir, "import.json")
	importData, err := json.MarshalIndent(testFacts, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal test facts: %v", err)
	}

	err = os.WriteFile(importFile, importData, 0644)
	if err != nil {
		t.Fatalf("Failed to write import file: %v", err)
	}

	// Test importing facts
	importCmd := &cobra.Command{
		Use:   "facts import",
		Short: "Import facts from JSON",
		RunE:  runFactsImport,
	}

	importTestCmd := NewTestCLICommand(importCmd)
	err = importTestCmd.Execute(importFile)
	if err != nil {
		t.Fatalf("Failed to execute facts import command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts import command executed successfully")

	// Verify facts were imported by querying them
	queryCmd := &cobra.Command{
		Use:   "facts query",
		Short: "Query stored facts",
		RunE:  runFactsQuery,
	}

	queryTestCmd := NewTestCLICommand(queryCmd)
	err = queryTestCmd.Execute("machine=test-machine-1")
	if err != nil {
		t.Fatalf("Failed to query imported facts: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts query command for imported facts executed successfully")
}

func TestFactsValidateCommand(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test validating facts
	validateCmd := &cobra.Command{
		Use:   "facts validate",
		Short: "Validate stored facts",
		RunE:  runFactsValidate,
	}

	validateTestCmd := NewTestCLICommand(validateCmd)
	err := validateTestCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute facts validate command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts validate command executed successfully")
}

func TestFactsCacheCommands(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test cache clear command
	clearCmd := &cobra.Command{
		Use:   "facts cache-clear",
		Short: "Clear fact cache",
		RunE:  runFactsCacheClear,
	}

	clearTestCmd := NewTestCLICommand(clearCmd)
	err := clearTestCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute facts cache-clear command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts cache-clear command executed successfully")

	// Test cache expired command
	expiredCmd := &cobra.Command{
		Use:   "facts cache-expired",
		Short: "Clear expired facts from cache",
		RunE:  runFactsCacheExpired,
	}

	expiredTestCmd := NewTestCLICommand(expiredCmd)
	err = expiredTestCmd.Execute()
	if err != nil {
		t.Fatalf("Failed to execute facts cache-expired command: %v", err)
	}

	// Since the command uses fmt.Println directly, we can't capture output easily
	// Instead, we verify that the command executed successfully without errors
	t.Log("Facts cache-expired command executed successfully")
}

func TestParseQueryExpression(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		expectError bool
		checkFunc   func(*testing.T, *facts.FactQuery)
	}{
		{
			name:       "simple machine query",
			expression: "machine=test-server",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if query.MachineName != "test-server" {
					t.Errorf("Expected machine name 'test-server', got '%s'", query.MachineName)
				}
			},
		},
		{
			name:       "simple OS query",
			expression: "os=Linux",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if query.OS != "Linux" {
					t.Errorf("Expected OS 'Linux', got '%s'", query.OS)
				}
			},
		},
		{
			name:       "simple project query",
			expression: "project=my-project",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if query.ProjectName != "my-project" {
					t.Errorf("Expected project name 'my-project', got '%s'", query.ProjectName)
				}
			},
		},
		{
			name:       "tag query with value",
			expression: "tag=environment:production",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if query.Tags["environment"] != "production" {
					t.Errorf("Expected tag environment=production, got %v", query.Tags)
				}
			},
		},
		{
			name:       "tag query without value",
			expression: "tag=environment",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if _, exists := query.Tags["environment"]; !exists {
					t.Errorf("Expected tag 'environment' to exist")
				}
			},
		},
		{
			name:       "multiple conditions",
			expression: "machine=test-server,os=Linux,tag=environment:production",
			checkFunc: func(t *testing.T, query *facts.FactQuery) {
				if query.MachineName != "test-server" {
					t.Errorf("Expected machine name 'test-server', got '%s'", query.MachineName)
				}
				if query.OS != "Linux" {
					t.Errorf("Expected OS 'Linux', got '%s'", query.OS)
				}
				if query.Tags["environment"] != "production" {
					t.Errorf("Expected tag environment=production, got %v", query.Tags)
				}
			},
		},
		{
			name:        "invalid query format",
			expression:  "invalid-format",
			expectError: true,
		},
		{
			name:        "unknown field",
			expression:  "unknown=value",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := parseQueryExpression(tt.expression)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, query)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
