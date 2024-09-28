package identity

import (
	"testing"
)

func TestIsSecure(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"Password1!", true},  // Valid password
		{"password", false},   // Missing uppercase, number, and special character
		{"PASSWORD", false},   // Missing lowercase, number, and special character
		{"Pass123", false},    // Missing special character
		{"P@ssw0rd", true},    // Valid password
		{"P@ssword", false},   // Missing number
		{"P4ssword", false},   // Missing special character
		{"P@ssw0rd123", true}, // Valid password with extra characters
		{"12345678", false},   // Missing uppercase, lowercase, and special character
		{"!@#$%^&*", false},   // Missing uppercase, lowercase, and number
		{"P@ssword$", false},  // Valid length, missing number
		{"Abcdefgh1", false},  // Missing special character
	}

	for _, test := range tests {
		result := IsSecure(test.password)
		if result != test.expected {
			t.Errorf("IsSecure(%q) = %v; want %v", test.password, result, test.expected)
		}
	}
}
