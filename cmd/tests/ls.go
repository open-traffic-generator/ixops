package tests

import (
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ixexec.ExecCmd("ls_tc")
		return nil
	},
}
