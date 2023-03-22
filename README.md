# aws-assume-role

## Installation

If you are using a Mac, you can use brew.

```zsh
$ brew tap naomichi-y/aws-assume-role
$ brew install naomichi-y/aws-assume-role/aws-assume-role
```

If the Go command is available.

```zsh
$ go install github.com/naomichi-y/aws-assume-role@latest
```

## Configuration

Define MFA enabled accounts in the profile.

`~/.aws/credentials`
```
[test]
aws_access_key_id = ***
aws_secret_access_key= ***
mfa_serial=(arn-of-the-mfa-device)
```

Temporary credentials can be created by running the aws-assume-role command and specifying an MFA profile and token.

```zsh
$ aws-assume-role
AWS profile [default]: test
Token code: ***
Access key ID: ***
Successfully updated test-assume profile. [~/.aws/credentials]
```

Temporary credentials are written to a credential file in the format `{PROFILE}-assume`.

```zsh
$ cat ~/.aws/credentials

[test-assume]
aws_access_key_id = ***
aws_secret_access_key = ***
aws_session_token = ***
```

You can execute AWS commands via Assume role.
```zsh
$ aws --profile test-assume sts get-caller-identity
{
    "UserId": "***,
    "Account": "***",
    "Arn": "***"
}
```
