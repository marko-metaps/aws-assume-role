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

const version = "1.0.2"

func checkCredentialFile() {
	dir, _ := os.UserHomeDir()
	path := dir + "/.aws/credentials"
	stat, err := os.Stat(path)

	if os.IsNotExist(err) {
		panic(err)
	}

	if stat.Mode().Perm() & 0200 != 0200 {
		panic(fmt.Errorf("No write permission. [~/.aws/credentials]"))
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
	var code string

	fmt.Print("Token code: ")
	fmt.Scanln(&code)

	return code
}

func getCredentials(p string, t string) credentials.Value {
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", p),
	}))

	sec, err := strconv.ParseInt(configureGetAlt(p, "duration_seconds", "3600"), 10, 64)

	if err != nil {
		panic(err)
	}

	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(sec),
		SerialNumber:    aws.String(configureGet(p, "mfa_serial")),
		TokenCode:       aws.String(t),
	}
	svc := sts.New(sess)
	token, err := svc.GetSessionToken(input)

	if err != nil {
		panic(err)
	}

	result := credentials.Value{
		AccessKeyID:     aws.StringValue(token.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(token.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(token.Credentials.SessionToken),
	}

	return result
}

func configureGet(p string, k string) string {
	result, err := exec.Command("aws", "--profile", p, "configure", "get", k).Output()

	if err != nil {
		panic(fmt.Errorf("%v is undefined in profile. [%v]\n", k, p))
	}

	return strings.ReplaceAll(string(result), "\n", "")
}

func configureGetAlt(p string, k string, v string) string {
	result := configureGet(p, k)

	if len(result) == 0 {
		return v
	}

	return result
}

func configureSet(p string, k string, v string) {
	_, err := exec.Command("aws", "--profile", p, "configure", "set", k, aws.StringValue(&v)).CombinedOutput()

	if err != nil {
		panic(err)
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
	profile := getProfile()

	code := getTokenCode()
	credentials := getCredentials(profile, code)

	fmt.Printf("Access key ID: %v\n", credentials.AccessKeyID)
	assume := profile + "-assume"

	configureSet(assume, "aws_access_key_id", credentials.AccessKeyID)
	configureSet(assume, "aws_secret_access_key", credentials.SecretAccessKey)
	configureSet(assume, "aws_session_token", credentials.SessionToken)

	fmt.Printf("Successfully updated %v profile. [~/.aws/credentials]\n", assume)
}
