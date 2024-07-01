package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/naomichi-y/aws-assume-role/pkg/config"
	creds "github.com/naomichi-y/aws-assume-role/pkg/credentials"
	"github.com/naomichi-y/aws-assume-role/pkg/filesystem"
	"github.com/naomichi-y/aws-assume-role/pkg/utils"
)

const version = "1.0.3"

func main() {
	var v bool

	flag.BoolVar(&v, "version", false, "Show version")
	flag.Parse()

	if v {
		fmt.Println(version)
		return
	}

	filesystem.CheckCredentialFile(filesystem.RealFileSystem{})

	profile := utils.GetProfile(os.Stdin)

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", profile),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create session: %s\n", err)
		os.Exit(1)
	}

	stsService := sts.New(sess)

	tokenCode := utils.GetTokenCode(os.Stdin)

	cmdRunner := config.RealCmdRunner{}
	credentials, err := creds.GetCredentials(stsService, cmdRunner, profile, tokenCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error retrieving credentials: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Access key ID: %v\n", credentials.AccessKeyID)
	assume := profile + "-assume"

	if err := config.ConfigureSet(cmdRunner, assume, "aws_access_key_id", credentials.AccessKeyID); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set aws_access_key_id: %s\n", err)
		os.Exit(1)
	}
	if err := config.ConfigureSet(cmdRunner, assume, "aws_secret_access_key", credentials.SecretAccessKey); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set aws_secret_access_key: %s\n", err)
		os.Exit(1)
	}
	if err := config.ConfigureSet(cmdRunner, assume, "aws_session_token", credentials.SessionToken); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set aws_session_token: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated %v profile. [~/.aws/credentials]\n", assume)
}
