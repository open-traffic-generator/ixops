package utils

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ExecCmd(cmd string, args ...string) (string, error) {
	var shellOutput bytes.Buffer
	var shellError bytes.Buffer
	errorString := ""
	command := exec.Command(cmd, args...)
	log.Printf("Executing: %v\n", command.Args)
	command.Stdout = &shellOutput
	command.Stderr = &shellError
	if err := command.Start(); err != nil {
		log.Printf("failed to start command execution: %v\n", err)
		return "", fmt.Errorf(fmt.Sprintf("failed to start command execution: %v", err))
	}

	err := command.Wait()
	log.Printf("Output: %v\n", shellOutput.String())
	if shellError.String() != "" {
		errorString += fmt.Sprintf("failed to wait for command to be executed: %s\n%v", shellError.String(), err)
	}

	log.Printf("Error: %v\n", errorString)

	if errorString != "" {
		return "", fmt.Errorf(errorString)
	}
	return shellOutput.String(), nil
}

func GetCommonHome() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get .ixops home: %v\n", err)
		return "", fmt.Errorf(fmt.Sprintf("failed to get .ixops home: %v\n", err))
	}
	return filepath.Join(userHome, ".ixops"), nil
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
