package utils

import (
	"bufio"
	"fmt"
	"io"
)

func GetProfile(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	fmt.Print("AWS profile [default]: ")
	scanner.Scan()
	profile := scanner.Text()

	if profile == "" {
		profile = "default"
	}

	return profile
}

func GetTokenCode(in io.Reader) (string, error) {
	fmt.Print("Token code: ")
	scanner := bufio.NewScanner(in)
	if scanner.Scan() {
		code := scanner.Text()
		if code == "" {
			return "", fmt.Errorf("no token code provided")
		}
		return code, nil
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading token code: %w", err)
	}
	return "", fmt.Errorf("no token code provided")
}
