package ixexec

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// to print the processed information when stdout gets a new line
func ExecCmd(commands string) {
	cmd := exec.Command("/home/ashukuma/athena/scripts/ixops/ixops.sh", strings.Fields(commands)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not start execution: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalf("Could not finish execution: %v", err)
	}
}
