package main

import (
	"bufio"
	"fmt"
	"io"
)

func getProfile(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	fmt.Print("AWS profile [default]: ")
	scanner.Scan()
	profile := scanner.Text()

	if profile == "" {
		profile = "default"
	}

	return profile
}

func getTokenCode(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	fmt.Print("Token code: ")
	scanner.Scan()
	return scanner.Text()
}
