package coordinator

import (
	"testing"

	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/types/logging"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorSecretsIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorSecretsIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestEncryptData(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorSecretsIntegration(nil, logger)

	data := []byte("test data")
	recipients := []string{"age1test"}

	result, err := integration.EncryptData(data, recipients)
	// Should return error when crypto manager is nil
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "crypto manager not available")
}

func TestDecryptData(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorSecretsIntegration(nil, logger)

	data := []byte("encrypted data")

	result, err := integration.DecryptData(data)
	// Should return error when crypto manager is nil
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "crypto manager not available")
}

func TestValidateEncryption(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorSecretsIntegration(nil, logger)

	// Test with valid encrypted data (at least 100 bytes)
	encryptedData := make([]byte, 150)
	for i := range encryptedData {
		encryptedData[i] = byte(i % 256)
	}

	err := integration.ValidateEncryption(encryptedData)
	// Should return error when crypto manager is nil
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "crypto manager not available")
}

func TestSecretsIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorSecretsIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 10
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			data := []byte("test data")
			recipients := []string{"age1test"}
			_, err := integration.EncryptData(data, recipients)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "crypto manager not available")
	}
}
