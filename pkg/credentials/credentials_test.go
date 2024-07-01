package creds

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/naomichi-y/aws-assume-role/testutils"
)

func TestGetCredentials(t *testing.T) {
	mockSTS := &testutils.MockSTSAPI{
		Resp: sts.GetSessionTokenOutput{
			Credentials: &sts.Credentials{
				AccessKeyId:     aws.String("AKIA1234567890"),
				SecretAccessKey: aws.String("secret"),
				SessionToken:    aws.String("token"),
			},
		},
		Err: nil,
	}

	mockRunner := &testutils.MockCmdRunner{
		MockOutput: map[string]string{
			"aws --profile default configure get mfa_serial":       "arn:aws:iam::123456789012:mfa/user",
			"aws --profile default configure get duration_seconds": "3600",
		},
		MockError: nil,
	}

	creds, err := GetCredentials(mockSTS, mockRunner, "default", "123456")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if creds.AccessKeyID != "AKIA1234567890" || creds.SecretAccessKey != "secret" || creds.SessionToken != "token" {
		t.Errorf("Unexpected credentials values: %+v", creds)
	}
}
