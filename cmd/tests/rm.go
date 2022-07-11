package tests

import (
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ixexec.ExecCmd("rm_tc")
		return nil
	},
}
