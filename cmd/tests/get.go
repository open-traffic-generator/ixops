package tests

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		// ixexec.ExecCmd("new_tc")
		return nil
	},
}
