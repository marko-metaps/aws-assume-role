package main

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/sts"
)

type STSAPI interface {
	GetSessionToken(input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error)
}

func getCredentials(stsService STSAPI, runner CmdRunner, profile string, tokenCode string) (credentials.Value, error) {
	durationSecondsStr := configureGetAlt(runner, profile, "duration_seconds", "3600")
	durationSeconds, err := strconv.ParseInt(durationSecondsStr, 10, 64)
	if err != nil {
		return credentials.Value{}, err
	}

	mfaSerial := configureGet(runner, profile, "mfa_serial")
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(durationSeconds),
		SerialNumber:    aws.String(mfaSerial),
		TokenCode:       aws.String(tokenCode),
	}

	token, err := stsService.GetSessionToken(input)
	if err != nil {
		return credentials.Value{}, err
	}

	return credentials.Value{
		AccessKeyID:     aws.StringValue(token.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(token.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(token.Credentials.SessionToken),
	}, nil
}
