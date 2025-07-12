package keygen

import (
	"strings"
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	// Test that password is generated with correct length
	password := GeneratePassword()
	if len(password) != 25 {
		t.Errorf("expected password length 25, got %d", len(password))
	}

	// Test that password contains only valid characters
	const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, char := range password {
		if !strings.ContainsRune(validChars, char) {
			t.Errorf("password contains invalid character: %c", char)
		}
	}

	// Test that multiple passwords are different (randomness)
	password1 := GeneratePassword()
	password2 := GeneratePassword()
	if password1 == password2 {
		t.Error("generated passwords should be different")
	}
}

func TestGeneratePasswordWithLength(t *testing.T) {
	testCases := []struct {
		name     string
		length   int
		expected int
	}{
		{"zero length", 0, 25},      // should default to 25
		{"negative length", -5, 25}, // should default to 25
		{"custom length 10", 10, 10},
		{"custom length 50", 50, 50},
		{"custom length 100", 100, 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			password := GeneratePasswordWithLength(tc.length)
			if len(password) != tc.expected {
				t.Errorf("expected password length %d, got %d", tc.expected, len(password))
			}

			// Test that password contains only valid characters
			const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			for _, char := range password {
				if !strings.ContainsRune(validChars, char) {
					t.Errorf("password contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestGeneratePasswordWithLength_Uniqueness(t *testing.T) {
	// Test that multiple passwords with same length are different
	password1 := GeneratePasswordWithLength(20)
	password2 := GeneratePasswordWithLength(20)
	if password1 == password2 {
		t.Error("generated passwords should be different")
	}
}
