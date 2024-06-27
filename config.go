package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type CmdRunner interface {
	RunCommand(name string, args ...string) ([]byte, error)
}

type RealCmdRunner struct{}

func (r RealCmdRunner) RunCommand(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).Output()
}

func configureGetRaw(runner CmdRunner, profile string, key string) string {
	result, _ := runner.RunCommand("aws", "--profile", profile, "configure", "get", key)
	return strings.TrimSpace(string(result))
}

func configureGet(runner CmdRunner, profile string, key string) string {
	result := configureGetRaw(runner, profile, key)
	if len(result) == 0 {
		panic(fmt.Errorf("%v is undefined in profile [%v]", key, profile))
	}
	return result
}

func configureGetAlt(runner CmdRunner, profile string, key string, defaultValue string) string {
	result := configureGetRaw(runner, profile, key)
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func configureSet(runner CmdRunner, profile, key, value string) error {
	_, err := runner.RunCommand("aws", "--profile", profile, "configure", "set", key, value)
	if err != nil {
		return fmt.Errorf("failed to execute 'aws --profile %s configure set %s %s': %v", profile, key, value, err)
	}
	return nil
}
