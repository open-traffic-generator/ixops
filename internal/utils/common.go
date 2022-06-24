package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

func WaitFor(fn func() (bool, error), opts *WaitForOpts) error {
	if opts == nil {
		opts = &WaitForOpts{
			Condition: "condition to be true",
		}
	}

	if opts.Interval == 0 {
		opts.Interval = 500 * time.Millisecond
	}
	if opts.Timeout == 0 {
		opts.Timeout = 300 * time.Second
	}

	start := time.Now()
	log.Printf("Waiting for %s ...\n", opts.Condition)

	for {
		done, err := fn()
		if err != nil {
			return fmt.Errorf("error waiting for %s: %v", opts.Condition, err)
		}
		if done {
			log.Printf("Done waiting for %s\n", opts.Condition)
			return nil
		}

		if time.Since(start) > opts.Timeout {
			return fmt.Errorf("timeout occurred while waiting for %s", opts.Condition)
		}
		time.Sleep(opts.Interval)
	}
}

func FileRelative(p string) (string, error) {
	bp, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return filepath.Dir(bp), nil
}

func ReadUrlAndWrite(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("%s file exists...", filePath)
		return true
	} else {
		log.Printf("%s file not exists...", filePath)
		return false
	}
}
