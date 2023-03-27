package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"os/exec"
	"fmt"
	"flag"
	"strings"
	"strconv"
)

const version = "1.0.1"

func checkCredentialFile() {
	h, _ := os.UserHomeDir()

	credentialsPath := h + "/.aws/credentials"
	_, err := os.Stat(credentialsPath)

	if os.IsNotExist(err) {
		fmt.Printf("%v File does not exist\n", credentialsPath)
		os.Exit(1)
	}
}

func getProfile() string {
	p := os.Getenv("AWS_PROFILE")

	if p == "" {
		fmt.Print("AWS profile [default]: ")
		fmt.Scanln(&p)

		if p == "" {
			p = "default"
		}
	}

	return p
}

func getTokenCode() string {
	var t string

	fmt.Print("Token code: ")
	fmt.Scanln(&t)

	return t
}

func getCredentials(p string, t string) credentials.Value {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", p),
	}))

	sec, err := strconv.ParseInt(configureGetAlt(p, "duration_seconds", "3600"), 10, 64)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(sec),
		SerialNumber:    aws.String(configureGet(p, "mfa_serial")),
		TokenCode:       aws.String(t),
	}
	svc := sts.New(sess)
	token, err := svc.GetSessionToken(input)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	v := credentials.Value{
		AccessKeyID:     aws.StringValue(token.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(token.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(token.Credentials.SessionToken),
	}

	return v
}

func configureGetRawData(p string, k string) string {
	result, _ := exec.Command("aws", "--profile", p, "configure", "get", k).Output()

	return strings.ReplaceAll(string(result), "\n", "")
}

func configureGet(p string, k string) string {
	result := configureGetRawData(p, k)

	if len(result) == 0 {
		fmt.Printf("%v is undefined in profile. [%v]\n", k, p)
		os.Exit(1)
	}

	return result
}

func configureGetAlt(p string, k string, alt string) string {
	result := configureGetRawData(p, k)

	if len(result) == 0 {
		return alt
	}

	return result
}

func configureSet(p string, k string, v string) {
	_, err := exec.Command("aws", "--profile", p, "configure", "set", k, aws.StringValue(&v)).CombinedOutput()

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func main() {
    var v bool

	flag.BoolVar(&v, "version", false, "Show version")
	flag.Parse()

	if v {
		fmt.Println(version)
		return
	}

	checkCredentialFile()
	p := getProfile()

	t := getTokenCode()
	credentials := getCredentials(p, t)

	fmt.Printf("Access key ID: %v\n", credentials.AccessKeyID)
	k := p + "-assume"

	configureSet(k, "aws_access_key_id", credentials.AccessKeyID)
	configureSet(k, "aws_secret_access_key", credentials.SecretAccessKey)
	configureSet(k, "aws_session_token", credentials.SessionToken)

	fmt.Printf("Successfully updated %v profile. [~/.aws/credentials]\n", k)
}
