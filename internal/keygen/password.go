package keygen

import (
	"crypto/rand"
	"math/big"
)

// GeneratePassword creates a password with ASCII letters and numbers
func GeneratePassword() string {
	return GeneratePasswordWithLength(DefaultPasswordLength)
}

// GeneratePasswordWithLength creates a password with the specified length
func GeneratePasswordWithLength(length int) string {
	// Handle edge cases
	if length <= 0 {
		length = DefaultPasswordLength
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)

	for i := range password {
		// Use crypto/rand for better randomness
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a simple method if crypto/rand fails
			password[i] = charset[i%len(charset)]
		} else {
			password[i] = charset[n.Int64()]
		}
	}
	return string(password)
}
