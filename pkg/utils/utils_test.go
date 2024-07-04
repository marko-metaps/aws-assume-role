package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestGetProfile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty input", "\n", "default"},
		{"non-default profile", "myprofile\n", "myprofile"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := bytes.NewBufferString(test.input)
			result := GetProfile(input)
			if result != test.expected {
				t.Errorf("expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestGetTokenCode(t *testing.T) {
	validInput := strings.NewReader("123456\n")
	code, err := GetTokenCode(validInput)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if code != "123456" {
		t.Errorf("Expected '123456', got '%s'", code)
	}

	// 空の入力
	emptyInput := strings.NewReader("\n")
	_, err = GetTokenCode(emptyInput)
	if err == nil {
		t.Errorf("Expected error for no input, got none")
	}
}
