package cmd_test

import (
	"os/exec"
	"testing"
)

func TestRootCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "help", args: []string{"help"}},
		{name: "bad", args: []string{"bad"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(ixOps, tt.args...)
			out, err := cmd.CombinedOutput()
			t.Logf("cmd output: %s\n", out)
			if tt.name == "bad" {
				if err == nil {
					t.Fatalf("cmd did not fail")
				}
			} else {
				if err != nil {
					t.Fatalf("cmd failed: %v\n", err)
				}
			}
		})
	}
}

func TestTopologyCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "topology", args: []string{"topology", "help"}},
		{name: "create", args: []string{"topology", "create", "help"}},
		{name: "delete", args: []string{"topology", "delete", "help"}},
		// {name: "bad", args: []string{"topology", "bad", "help"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(ixOps, tt.args...)
			out, err := cmd.CombinedOutput()
			t.Logf("cmd output: %s\n", out)
			if tt.name == "bad" {
				if err == nil {
					t.Fatalf("cmd did not fail")
				}
			} else {
				if err != nil {
					t.Fatalf("cmd failed: %v\n", err)
				}
			}
		})
	}
}

func TestClusterCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "cluster", args: []string{"cluster", "help"}},
		{name: "setup", args: []string{"cluster", "setup", "help"}},
		{name: "teardown", args: []string{"cluster", "teardown", "help"}},
		// {name: "bad", args: []string{"cluster", "bad", "help"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(ixOps, tt.args...)
			out, err := cmd.CombinedOutput()
			t.Logf("cmd output: %s\n", out)
			if tt.name == "bad" {
				if err == nil {
					t.Fatalf("cmd did not fail")
				}
			} else {
				if err != nil {
					t.Fatalf("cmd failed: %v\n", err)
				}
			}
		})
	}
}

func TestImagesCmd(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "images", args: []string{"images", "help"}},
		{name: "get", args: []string{"images", "get", "help"}},
		{name: "rm", args: []string{"images", "rm", "help"}},
		// {name: "bad", args: []string{"images", "bad", "help"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(ixOps, tt.args...)
			out, err := cmd.CombinedOutput()
			t.Logf("cmd output: %s\n", out)
			if tt.name == "bad" {
				if err == nil {
					t.Fatalf("cmd did not fail")
				}
			} else {
				if err != nil {
					t.Fatalf("cmd failed: %v\n", err)
				}
			}
		})
	}
}
