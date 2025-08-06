package coordinator

import (
	"testing"

	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/logging/types"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorCryptoIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorCryptoIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestEncryptData(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorCryptoIntegration(nil, logger)

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
	integration := NewCoordinatorCryptoIntegration(nil, logger)

	data := []byte("encrypted data")

	result, err := integration.DecryptData(data)
	// Should return error when crypto manager is nil
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "crypto manager not available")
}

func TestValidateEncryption(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorCryptoIntegration(nil, logger)

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

func TestGetCryptoStatus(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorCryptoIntegration(nil, logger)

	status := integration.GetCryptoStatus()
	assert.NotNil(t, status)
	assert.IsType(t, map[string]interface{}{}, status)
}

func TestCryptoIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorCryptoIntegration(nil, logger)

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
