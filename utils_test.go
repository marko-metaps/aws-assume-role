package main

import (
	"bytes"
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
			result := getProfile(input)
			if result != test.expected {
				t.Errorf("expected '%s', got '%s'", test.expected, result)
			}
		})
	}
}

func TestGetTokenCode(t *testing.T) {
	token := "123456"
	input := bytes.NewBufferString(token + "\n")
	result := getTokenCode(input)
	if result != token {
		t.Errorf("expected '%s', got '%s'", token, result)
	}
}
