package tests

import (
	"github.com/open-traffic-generator/ixops/pkg/ixexec"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ixexec.ExecCmd("run_tc " + args[0])
		return nil
	},
}
