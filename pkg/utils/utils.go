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

func GetTokenCode(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	fmt.Print("Token code: ")
	scanner.Scan()
	return scanner.Text()
}
