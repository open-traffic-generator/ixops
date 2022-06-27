package tests

import (
	"github.com/open-traffic-generator/ixops/internal/run_tests"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := run_tests.RunTests(args)
		if err != nil {
			return err
		}
		return nil
	},
}
