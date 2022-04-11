package diagnostics

import (
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect diagnostics",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
