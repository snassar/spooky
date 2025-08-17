package secrets

import (
	"context"
	"fmt"
	"strings"
	"testing"

	spookytypes "spooky/internal/types"
	spookytypesfacts "spooky/internal/types/facts"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypesvariables "spooky/internal/types/variables"

	"github.com/stretchr/testify/assert"
)

// Test constants
const (
	testDecryptedFact   = "decrypted_fact"
	testDecryptedCustom = "decrypted_custom"
)

// MockSecretsIntegration implements SecretsIntegration for testing
type MockSecretsIntegration struct {
	decryptFunc func(ctx context.Context, data []byte, identityPath string) ([]byte, error)
	encryptFunc func(ctx context.Context, data []byte, recipients []string) ([]byte, error)
}

func (m *MockSecretsIntegration) EncryptWithAge(_ context.Context, data []byte, recipients []string) ([]byte, error) {
	if m.encryptFunc != nil {
		return m.encryptFunc(context.Background(), data, recipients)
	}
	return []byte("age1encrypted"), nil
}

func (m *MockSecretsIntegration) DecryptWithAge(_ context.Context, data []byte, identityPath string) ([]byte, error) {
	if m.decryptFunc != nil {
		return m.decryptFunc(context.Background(), data, identityPath)
	}
	// Simple mock decryption - remove "age1encrypted_" prefix
	if strings.HasPrefix(string(data), "age1encrypted_") {
		return []byte(strings.TrimPrefix(string(data), "age1encrypted_")), nil
	}
	return data, nil
}

func (m *MockSecretsIntegration) EncryptWithPassphrase(_ context.Context, data []byte, _ string) ([]byte, error) {
	return []byte("encrypted_" + string(data)), nil
}

func (m *MockSecretsIntegration) DecryptWithPassphrase(_ context.Context, data []byte, _ string) ([]byte, error) {
	return []byte(strings.TrimPrefix(string(data), "encrypted_")), nil
}

func (m *MockSecretsIntegration) ValidateAgeKey(_ context.Context, _ string) error {
	return nil
}

func (m *MockSecretsIntegration) ListRecipients(_ context.Context, _ []byte) ([]string, error) {
	return nil, nil
}

func (m *MockSecretsIntegration) ValidateAgeEncryptedValue(_ context.Context, _ string) error {
	return nil
}

func (m *MockSecretsIntegration) LoadRecipients(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

func (m *MockSecretsIntegration) LoadIdentities(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

func (m *MockSecretsIntegration) EncryptHCLValues(ctx context.Context, data interface{}, recipients []string, dryRun bool) error {
	processor := NewHCLProcessor(&MockLogger{})
	return processor.EncryptHCLValues(ctx, data, m, recipients, dryRun)
}

func (m *MockSecretsIntegration) EncryptHCLValuesSensitive(ctx context.Context, data interface{}, recipients []string, dryRun bool, _ func(path []string, value interface{}) bool) error {
	processor := NewHCLProcessor(&MockLogger{})
	return processor.EncryptHCLValues(ctx, data, m, recipients, dryRun)
}

func (m *MockSecretsIntegration) EncryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, recipients []string, dryRun bool) error {
	processor := NewHCLProcessor(&MockLogger{})
	return processor.EncryptHCLValues(ctx, data, m, recipients, dryRun)
}

func (m *MockSecretsIntegration) DecryptHCLValues(ctx context.Context, data interface{}, identityPath string) error {
	processor := NewHCLProcessor(&MockLogger{})
	return processor.DecryptHCLValues(ctx, data, m, identityPath)
}

func (m *MockSecretsIntegration) DecryptHCLValuesWithJSONSupport(ctx context.Context, data interface{}, identityPath string) error {
	processor := NewHCLProcessor(&MockLogger{})
	return processor.DecryptHCLValues(ctx, data, m, identityPath)
}

// MockLogger implements Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, fields ...map[string]interface{}) {
	fmt.Printf("DEBUG: %s %+v\n", msg, fields)
}
func (m *MockLogger) Info(msg string, fields ...map[string]interface{}) {
	fmt.Printf("INFO: %s %+v\n", msg, fields)
}
func (m *MockLogger) Warn(msg string, fields ...map[string]interface{}) {
	fmt.Printf("WARN: %s %+v\n", msg, fields)
}
func (m *MockLogger) Error(msg string, err error, fields ...map[string]interface{}) {
	fmt.Printf("ERROR: %s %v %+v\n", msg, err, fields)
}
func (m *MockLogger) Fatal(msg string, err error, fields ...map[string]interface{}) {
	fmt.Printf("FATAL: %s %v %+v\n", msg, err, fields)
}
func (m *MockLogger) WithFields(_ map[string]interface{}) spookytypes.Logger { return m }
func (m *MockLogger) WithComponent(_ string) spookytypes.Logger              { return m }
func (m *MockLogger) WithOperation(_ string) spookytypes.Logger              { return m }
func (m *MockLogger) SetLevel(_ spookytypes.LogLevel)                        {}
func (m *MockLogger) GetLevel() spookytypes.LogLevel                         { return spookytypeslogging.LogLevelInfo }

func TestDecryptHCLValues_ReplacesDecryptVariables(t *testing.T) {
	// Create mock secrets integration
	secretsIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			// Mock decryption - return "decrypted_value" for any age1... string
			return []byte("decrypted_value"), nil
		},
	}

	// Create variables data structure (similar to what DecryptVariables processes)
	variables := map[string]*spookytypesvariables.Variable{
		"database_password": {
			Name:      "database_password",
			Type:      "string",
			Default:   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			Encrypted: true,
		},
		"api_key": {
			Name:      "api_key",
			Type:      "string",
			Default:   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			Encrypted: true,
		},
	}

	// Use the generic DecryptHCLValues instead of DecryptVariables
	ctx := context.Background()
	err := secretsIntegration.DecryptHCLValues(ctx, variables, "/path/to/identity")

	if err != nil {
		t.Fatalf("DecryptHCLValues failed: %v", err)
	}

	// Verify that the encrypted values were decrypted
	for name, variable := range variables {
		if variable.Encrypted {
			if strValue, ok := variable.Default.(string); ok {
				if strValue != "decrypted_value" {
					t.Errorf("Variable %s was not decrypted properly: got %s, want decrypted_value", name, strValue)
				}
			}
		}
	}
}

func TestDecryptHCLValues_ReplacesDecryptMachines(t *testing.T) {
	// Create mock secrets integration
	secretsIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			// Mock decryption - return "decrypted_secret" for any age1... string
			return []byte("decrypted_secret"), nil
		},
	}

	// Create machines data structure (similar to what DecryptMachines processes)
	machines := []spookytypes.Machine{
		{
			Hostname:   "web-server",
			Password:   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			Passphrase: "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
		},
		{
			Hostname:   "db-server",
			Password:   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			Passphrase: "",
		},
	}

	// Use the generic DecryptHCLValues instead of DecryptMachines
	ctx := context.Background()
	err := secretsIntegration.DecryptHCLValues(ctx, &machines, "/path/to/identity")

	if err != nil {
		t.Fatalf("DecryptHCLValues failed: %v", err)
	}

	// Verify that the encrypted values were decrypted
	for i, machine := range machines {
		if machine.Password != "" && machine.Password != "decrypted_secret" {
			t.Errorf("Machine %d password was not decrypted properly: got %s, want decrypted_secret", i, machine.Password)
		}
		if machine.Passphrase != "" && machine.Passphrase != "decrypted_secret" {
			t.Errorf("Machine %d passphrase was not decrypted properly: got %s, want decrypted_secret", i, machine.Passphrase)
		}
	}
}

func TestDecryptHCLValues_ReplacesDecryptFacts(t *testing.T) {
	// Create mock secrets integration
	secretsIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			// Mock decryption - return testDecryptedFact for any age1... string
			return []byte(testDecryptedFact), nil
		},
	}

	// Create facts data structure (similar to what DecryptFacts processes)
	facts := &spookytypesfacts.FactCollection{
		MachineID: "test-machine",
		Facts: &spookytypesfacts.Facts{
			Custom: map[string]interface{}{
				"secret_key": "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
				"nested_secrets": map[string]interface{}{
					"api_token": "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
				},
				"secret_array": []interface{}{
					"age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
					"plain_value",
				},
			},
		},
	}

	// Use the generic DecryptHCLValues instead of DecryptFacts
	ctx := context.Background()
	err := secretsIntegration.DecryptHCLValues(ctx, facts, "/path/to/identity")

	if err != nil {
		t.Fatalf("DecryptHCLValues failed: %v", err)
	}

	// Verify that the encrypted values were decrypted
	customFacts := facts.Facts.Custom
	if customFacts["secret_key"] != testDecryptedFact {
		t.Errorf("secret_key was not decrypted properly: got %v, want %s", customFacts["secret_key"], testDecryptedFact)
	}

	nestedSecrets := customFacts["nested_secrets"].(map[string]interface{})
	if nestedSecrets["api_token"] != testDecryptedFact {
		t.Errorf("nested api_token was not decrypted properly: got %v, want %s", nestedSecrets["api_token"], testDecryptedFact)
	}

	secretArray := customFacts["secret_array"].([]interface{})
	if secretArray[0] != testDecryptedFact {
		t.Errorf("secret_array[0] was not decrypted properly: got %v, want %s", secretArray[0], testDecryptedFact)
	}
	if secretArray[1] != "plain_value" {
		t.Errorf("secret_array[1] was modified when it shouldn't be: got %v, want plain_value", secretArray[1])
	}
}

func TestDecryptHCLValues_HandlesCustomHCLFiles(t *testing.T) {
	// Create mock secrets integration
	secretsIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			// Mock decryption - return testDecryptedCustom for any age1... string
			return []byte(testDecryptedCustom), nil
		},
	}

	// Create custom HCL data structure (similar to /etc/spooky/custom.hcl on remote machines)
	customHCL := map[string]interface{}{
		"environment": "production",
		"secret_config": map[string]interface{}{
			"database_url": "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			"api_secret":   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
		},
		"public_config": map[string]interface{}{
			"log_level": "info",
			"port":      8080,
		},
		"secret_list": []interface{}{
			"age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			"age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
		},
	}

	// Use the generic DecryptHCLValues for custom HCL files
	ctx := context.Background()
	err := secretsIntegration.DecryptHCLValues(ctx, customHCL, "/path/to/identity")

	if err != nil {
		t.Fatalf("DecryptHCLValues failed: %v", err)
	}

	// Verify that only encrypted values were decrypted
	if customHCL["environment"] != "production" {
		t.Errorf("environment was modified when it shouldn't be: got %v, want production", customHCL["environment"])
	}

	secretConfig := customHCL["secret_config"].(map[string]interface{})
	if secretConfig["database_url"] != testDecryptedCustom {
		t.Errorf("database_url was not decrypted properly: got %v, want %s", secretConfig["database_url"], testDecryptedCustom)
	}
	if secretConfig["api_secret"] != testDecryptedCustom {
		t.Errorf("api_secret was not decrypted properly: got %v, want %s", secretConfig["api_secret"], testDecryptedCustom)
	}

	publicConfig := customHCL["public_config"].(map[string]interface{})
	if publicConfig["log_level"] != "info" {
		t.Errorf("log_level was modified when it shouldn't be: got %v, want info", publicConfig["log_level"])
	}
	if publicConfig["port"] != 8080 {
		t.Errorf("port was modified when it shouldn't be: got %v, want 8080", publicConfig["port"])
	}

	secretList := customHCL["secret_list"].([]interface{})
	if secretList[0] != testDecryptedCustom {
		t.Errorf("secret_list[0] was not decrypted properly: got %v, want %s", secretList[0], testDecryptedCustom)
	}
	if secretList[1] != testDecryptedCustom {
		t.Errorf("secret_list[1] was not decrypted properly: got %v, want %s", secretList[1], testDecryptedCustom)
	}
}

func TestEncryptHCLValues_EncryptsStringValues(t *testing.T) {
	// Create test data with unencrypted strings
	testData := map[string]interface{}{
		"password": "secret123",
		"api_key":  "key456",
		"nested": map[string]interface{}{
			"token": "nested_secret",
		},
	}

	mockIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			return []byte("decrypted"), nil
		},
		encryptFunc: func(_ context.Context, _ []byte, _ []string) ([]byte, error) {
			return []byte("age1encrypted"), nil
		},
	}

	processor := NewHCLProcessor(&MockLogger{})

	// Test dry run first
	err := processor.EncryptHCLValues(context.Background(), testData, mockIntegration, []string{"recipient1"}, true)
	assert.NoError(t, err)

	// Test actual encryption
	err = processor.EncryptHCLValues(context.Background(), testData, mockIntegration, []string{"recipient1"}, false)
	assert.NoError(t, err)

	// Verify that string values were encrypted
	assert.Equal(t, "age1encrypted", testData["password"])
	assert.Equal(t, "age1encrypted", testData["api_key"])
	assert.Equal(t, "age1encrypted", testData["nested"].(map[string]interface{})["token"])
}

func TestDecryptHCLValuesWithJSONSupport_HandlesJSONSerializedValues(t *testing.T) {
	// Create mock secrets integration that returns JSON-encoded data
	secretsIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			// Mock decryption - return JSON-encoded complex data
			return []byte(`{"name":"test","value":42,"enabled":true}`), nil
		},
	}

	// Create variables with JSON-serialized encrypted values
	variables := map[string]*spookytypesvariables.Variable{
		"complex_config": {
			Name:      "complex_config",
			Type:      "object",
			Default:   "age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p",
			Encrypted: true,
		},
	}

	t.Logf("Before decryption: Default = %v (type: %T)", variables["complex_config"].Default, variables["complex_config"].Default)

	// Use DecryptHCLValuesWithJSONSupport for JSON-serialized values
	ctx := context.Background()
	err := secretsIntegration.DecryptHCLValuesWithJSONSupport(ctx, variables, "/path/to/identity")

	if err != nil {
		t.Fatalf("DecryptHCLValuesWithJSONSupport failed: %v", err)
	}

	t.Logf("After decryption: Default = %v (type: %T)", variables["complex_config"].Default, variables["complex_config"].Default)

	// Verify that the JSON was deserialized properly
	complexConfig := variables["complex_config"].Default.(map[string]interface{})
	if complexConfig["name"] != "test" {
		t.Errorf("complex_config.name was not deserialized properly: got %v, want test", complexConfig["name"])
	}
	if complexConfig["value"] != float64(42) {
		t.Errorf("complex_config.value was not deserialized properly: got %v, want 42", complexConfig["value"])
	}
	if complexConfig["enabled"] != true {
		t.Errorf("complex_config.enabled was not deserialized properly: got %v, want true", complexConfig["enabled"])
	}
}

func TestEncryptHCLValues_HandlesComplexVariableValues(t *testing.T) {
	// Test complex variable values that include maps, objects, and nested structures
	variables := map[string]*spookytypesvariables.Variable{
		"simple_string": {
			Name:      "simple_string",
			Default:   "secret123",
			Type:      spookytypesvariables.VariableTypeString,
			Encrypted: true, // This variable should be encrypted
		},
		"database_config": {
			Name:      "database_config",
			Type:      spookytypesvariables.VariableTypeMap,
			Encrypted: true, // This variable should be encrypted
			Default: map[string]interface{}{
				"host":     "localhost",
				"port":     5432,
				"username": "admin",
				"password": "secret_password",
				"ssl":      true,
			},
		},
		"api_credentials": {
			Name:      "api_credentials",
			Type:      spookytypesvariables.VariableTypeObject,
			Encrypted: false, // This variable should NOT be encrypted
			Default: map[string]interface{}{
				"endpoint": "https://api.example.com",
				"auth": map[string]interface{}{
					"type":    "bearer",
					"token":   "secret_token_here", // Should NOT be encrypted (variable.Encrypted = false)
					"expires": "2024-12-31",
				},
				"headers": map[string]interface{}{
					"user-agent": "spooky-client",
					"x-api-key":  "secret_api_key", // Should NOT be encrypted (variable.Encrypted = false)
				},
			},
		},
		"server_list": {
			Name:      "server_list",
			Type:      spookytypesvariables.VariableTypeList,
			Encrypted: true, // This variable should be encrypted
			Default: []interface{}{
				"server1.example.com",
				"server2.example.com",
				"secret_internal_server",
			},
		},
		"mixed_config": {
			Name:      "mixed_config",
			Type:      spookytypesvariables.VariableTypeObject,
			Encrypted: false, // This variable should NOT be encrypted
			Default: map[string]interface{}{
				"public_info": map[string]interface{}{
					"name":    "my-app",
					"version": "1.0.0",
				},
				"private_info": map[string]interface{}{
					"secret_key": "very_secret_key", // Should NOT be encrypted (variable.Encrypted = false)
					"tokens": []interface{}{
						"token1",
						"secret_token2", // Should NOT be encrypted (variable.Encrypted = false)
						"token3",
					},
				},
			},
		},
	}

	mockIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			return []byte("decrypted_value"), nil
		},
		encryptFunc: func(_ context.Context, data []byte, _ []string) ([]byte, error) {
			return []byte("age1encrypted_" + string(data)), nil
		},
	}

	processor := NewHCLProcessor(&MockLogger{})

	// Test explicit encryption based on encrypted flag
	err := processor.EncryptHCLValues(context.Background(), variables, mockIntegration, []string{"recipient1"}, false)
	assert.NoError(t, err)

	// Verify that variables with encrypted = true were encrypted
	assert.Equal(t, "age1encrypted_secret123", variables["simple_string"].Default)

	// Verify map values were encrypted (only for variables with encrypted = true)
	dbConfig := variables["database_config"].Default.(map[string]interface{})
	assert.Equal(t, "age1encrypted_localhost", dbConfig["host"])
	assert.Equal(t, 5432, dbConfig["port"]) // Numbers are not encrypted
	assert.Equal(t, "age1encrypted_admin", dbConfig["username"])
	assert.Equal(t, "age1encrypted_secret_password", dbConfig["password"])
	assert.Equal(t, true, dbConfig["ssl"]) // Booleans are not encrypted

	// Verify variables with encrypted = false were NOT encrypted
	apiCreds := variables["api_credentials"].Default.(map[string]interface{})
	auth := apiCreds["auth"].(map[string]interface{})
	headers := apiCreds["headers"].(map[string]interface{})
	assert.Equal(t, "secret_token_here", auth["token"])              // NOT encrypted
	assert.Equal(t, "secret_api_key", headers["x-api-key"])          // NOT encrypted
	assert.Equal(t, "https://api.example.com", apiCreds["endpoint"]) // NOT encrypted

	// Verify list values were encrypted (only for variables with encrypted = true)
	serverList := variables["server_list"].Default.([]interface{})
	assert.Equal(t, "age1encrypted_server1.example.com", serverList[0])
	assert.Equal(t, "age1encrypted_server2.example.com", serverList[1])
	assert.Equal(t, "age1encrypted_secret_internal_server", serverList[2])

	// Verify variables with encrypted = false were NOT encrypted
	mixedConfig := variables["mixed_config"].Default.(map[string]interface{})
	privateInfo := mixedConfig["private_info"].(map[string]interface{})
	tokens := privateInfo["tokens"].([]interface{})
	assert.Equal(t, "very_secret_key", privateInfo["secret_key"]) // NOT encrypted
	assert.Equal(t, "secret_token2", tokens[1])                   // NOT encrypted
	assert.Equal(t, "token1", tokens[0])                          // NOT encrypted

	// Test decryption
	err = processor.DecryptHCLValues(context.Background(), variables, mockIntegration, "test-identity")
	assert.NoError(t, err)

	// Verify that encrypted values were decrypted
	assert.Equal(t, "decrypted_value", variables["simple_string"].Default)
	assert.Equal(t, "decrypted_value", dbConfig["host"])
	assert.Equal(t, 5432, dbConfig["port"]) // Numbers are not encrypted/decrypted
	assert.Equal(t, "decrypted_value", dbConfig["username"])
	assert.Equal(t, "decrypted_value", dbConfig["password"])
	assert.Equal(t, true, dbConfig["ssl"]) // Booleans are not encrypted/decrypted
	assert.Equal(t, "decrypted_value", serverList[0])
	assert.Equal(t, "decrypted_value", serverList[1])
	assert.Equal(t, "decrypted_value", serverList[2])

	// Verify that non-encrypted values remain unchanged
	assert.Equal(t, "secret_token_here", auth["token"])
	assert.Equal(t, "secret_api_key", headers["x-api-key"])
	assert.Equal(t, "very_secret_key", privateInfo["secret_key"])
	assert.Equal(t, "secret_token2", tokens[1])
}

func TestEncryptHCLValuesWithJSONSupport_EncryptsComplexValues(t *testing.T) {
	// Test complex values that need JSON serialization before encryption
	data := map[string]interface{}{
		"simple_string": "secret123",
		"complex_object": map[string]interface{}{
			"nested_map": map[string]interface{}{
				"key1": "value1",
				"key2": "secret_value",
			},
			"nested_array": []interface{}{
				"item1",
				"secret_item",
				"item3",
			},
		},
		"mixed_types": []interface{}{
			"string_value",
			42,
			true,
			map[string]interface{}{
				"nested": "secret_nested",
			},
		},
	}

	mockIntegration := &MockSecretsIntegration{
		decryptFunc: func(_ context.Context, _ []byte, _ string) ([]byte, error) {
			return []byte("decrypted_value"), nil
		},
		encryptFunc: func(_ context.Context, data []byte, _ []string) ([]byte, error) {
			return []byte("age1encrypted_" + string(data)), nil
		},
	}

	processor := NewHCLProcessor(&MockLogger{})

	// Test encryption with JSON support
	err := processor.EncryptHCLValues(context.Background(), data, mockIntegration, []string{"recipient1"}, false)
	assert.NoError(t, err)

	// Verify that string values were encrypted
	assert.Equal(t, "age1encrypted_secret123", data["simple_string"])

	// Verify that complex objects were JSON-serialized and encrypted
	complexObj := data["complex_object"].(string)
	assert.True(t, strings.HasPrefix(complexObj, "age1encrypted_"))

	// Verify that mixed types were JSON-serialized and encrypted
	mixedTypes := data["mixed_types"].(string)
	assert.True(t, strings.HasPrefix(mixedTypes, "age1encrypted_"))

	// Test decryption with JSON support
	err = processor.DecryptHCLValues(context.Background(), data, mockIntegration, "test-identity")
	assert.NoError(t, err)

	// Verify that encrypted values were decrypted
	assert.Equal(t, "decrypted_value", data["simple_string"])
	assert.Equal(t, "decrypted_value", data["complex_object"])
	assert.Equal(t, "decrypted_value", data["mixed_types"])
}

func TestEncryptHCLValues_Debug(t *testing.T) {
	// Simple test to debug encryption
	testData := map[string]interface{}{
		"test": "value",
	}

	mockIntegration := &MockSecretsIntegration{
		encryptFunc: func(_ context.Context, data []byte, _ []string) ([]byte, error) {
			t.Logf("EncryptWithAge called with data: %s", string(data))
			return []byte("age1encrypted"), nil
		},
	}

	// Create a logger that actually logs
	logger := &MockLogger{}
	processor := NewHCLProcessor(logger)

	t.Logf("Before encryption: %+v", testData)

	// Test actual encryption
	err := processor.EncryptHCLValues(context.Background(), testData, mockIntegration, []string{"recipient1"}, false)
	assert.NoError(t, err)

	t.Logf("After encryption: %+v", testData)
	assert.Equal(t, "age1encrypted", testData["test"])
}
