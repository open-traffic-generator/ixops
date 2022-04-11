package cmd_test

import (
	"os/exec"
	"path"
	"testing"
)

var (
	rootDir = "../.."
	ixOps   = path.Join(".", "ixops")
)

// TestMain is the first thing that's executed upon running `go test ...`
func TestMain(t *testing.T) {
	cmd := exec.Command("go", "build")
	cmd.Dir = rootDir
	out, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("cmd.Run() failed: %v\n: out%s\n", err, out)
	}
}
