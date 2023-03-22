package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"os/exec"
	"fmt"
	"strings"
)

func main() {
	homeDir, _ := os.UserHomeDir()

	credentialsPath := homeDir + "/.aws/credentials"
	_, err := os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		fmt.Printf("%v File does not exist\n", credentialsPath)
		os.Exit(1)
	}

	aws_profile := os.Getenv("AWS_PROFILE")

	if aws_profile == "" {
		fmt.Print("AWS profile [default]: ")
		fmt.Scanln(&aws_profile)

		if aws_profile == "" {
			aws_profile = "default"
		}
	}

	mfa_serial, _ := exec.Command("aws", "--profile", aws_profile, "configure", "get", "mfa_serial").Output()

	if len(mfa_serial) == 0 {
		fmt.Printf("`mfa_serial` is undefined in profile. [%v]\n", aws_profile)
		os.Exit(1)
	}

	var token_code string

	fmt.Print("Token code: ")
	fmt.Scanln(&token_code)

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", aws_profile),
	}))
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(3600),
		SerialNumber:    aws.String(strings.ReplaceAll(string(mfa_serial), "\n", "")),
		TokenCode:       aws.String(token_code),
	}
	svc := sts.New(sess)
	result, err := svc.GetSessionToken(input)

	if (err != nil) {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Access key ID: %v\n", aws.StringValue(result.Credentials.AccessKeyId))
	var assume_new_profile = aws_profile + "-assume"

	_, err = exec.Command("aws", "--profile", assume_new_profile, "configure", "set", "aws_access_key_id", aws.StringValue(result.Credentials.AccessKeyId)).CombinedOutput()
	_, err = exec.Command("aws", "--profile", assume_new_profile, "configure", "set", "aws_secret_access_key", aws.StringValue(result.Credentials.SecretAccessKey)).CombinedOutput()
	_, err = exec.Command("aws", "--profile", assume_new_profile, "configure", "set", "aws_session_token", aws.StringValue(result.Credentials.SessionToken)).CombinedOutput()

	if (err != nil) {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated %v profile. [~/.aws/credentials]\n", assume_new_profile)
}
