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

func checkCredentialFile() {
	homeDir, _ := os.UserHomeDir()

	credentialsPath := homeDir + "/.aws/credentials"
	_, err := os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		fmt.Printf("%v File does not exist\n", credentialsPath)
		os.Exit(1)
	}
}

func getProfile() string {
	profile := os.Getenv("AWS_PROFILE")

	if profile == "" {
		fmt.Print("AWS profile [default]: ")
		fmt.Scanln(&profile)

		if profile == "" {
			profile = "default"
		}
	}

	return profile
}

func getTokenCode() string {
	var tokenCode string

	fmt.Print("Token code: ")
	fmt.Scanln(&tokenCode)

	return tokenCode
}

func getCredentials(profile string, mfaSerial string, tokenCode string) credentials.Value {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", profile),
	}))
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(3600),
		SerialNumber:    aws.String(mfaSerial),
		TokenCode:       aws.String(tokenCode),
	}
	svc := sts.New(sess)
	token, err := svc.GetSessionToken(input)

	if (err != nil) {
		fmt.Print(err)
		os.Exit(1)
	}

	value := credentials.Value{
		AccessKeyID:     aws.StringValue(token.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(token.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(token.Credentials.SessionToken),
	}
	return value
}

func configureGet(profile string, key string) string {
	result, _ := exec.Command("aws", "--profile", profile, "configure", "get", key).Output()

	if len(result) == 0 {
		fmt.Printf("%v is undefined in profile. [%v]\n", key, profile)
		os.Exit(1)
	}

	return strings.ReplaceAll(string(result), "\n", "")
}

func configureSet(profile string, key string, value string) {
	_, err := exec.Command("aws", "--profile", profile, "configure", "set", key, aws.StringValue(&value)).CombinedOutput()

	if (err != nil) {
		fmt.Print(err)
		os.Exit(1)
	}
}

func main() {
	checkCredentialFile()
	profile := getProfile()

	mfaSerial := configureGet(profile, "mfa_serial")
	tokenCode := getTokenCode()
	credentials := getCredentials(profile, mfaSerial, tokenCode)

	fmt.Printf("Access key ID: %v\n", credentials.AccessKeyID)
	var newProfile = profile + "-assume"

	configureSet(newProfile, "aws_access_key_id", credentials.AccessKeyID)
	configureSet(newProfile, "aws_secret_access_key", credentials.SecretAccessKey)
	configureSet(newProfile, "aws_session_token", credentials.SessionToken)

	fmt.Printf("Successfully updated %v profile. [~/.aws/credentials]\n", newProfile)
}
