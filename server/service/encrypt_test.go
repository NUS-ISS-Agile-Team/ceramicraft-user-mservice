package service

import (
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	password := "testPassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		expectError    bool
	}{
		{"Valid password", hashedPassword, password, false},
		{"Invalid password", hashedPassword, "wrongPassword", true},
		{"Empty password", hashedPassword, "", true},
		{"Empty hashed password", "", password, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("VerifyPassword() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{"Valid password", "testPassword", false},
		{"Empty password", "", false},
		{"Long password", "aVeryLongPassword1234567890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)
			if (err != nil) != tt.expectError {
				t.Errorf("HashPassword() error = %v, expectError %v", err, tt.expectError)
			}
			if !tt.expectError && hashedPassword == "" {
				t.Error("HashPassword() returned an empty hashed password")
			}
		})
	}
}
