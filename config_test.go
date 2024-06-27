package main

import (
	"testing"

	"github.com/naomichi-y/aws-assume-role/testutils"
)

func TestConfigureGet(t *testing.T) {
	mockRunner := &testutils.MockCmdRunner{
		MockOutput: map[string]string{
			"aws --profile default configure get aws_access_key_id": "AKIA1234567890",
		},
		MockError: nil,
	}

	key := "aws_access_key_id"
	got := configureGet(mockRunner, "default", key)
	want := "AKIA1234567890"
	if got != want {
		t.Errorf("configureGet() = %v, want %v", got, want)
	}

	mockRunner.MockOutput = map[string]string{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for undefined key, did not panic")
		}
	}()
	configureGet(mockRunner, "default", "undefined_key")
}
