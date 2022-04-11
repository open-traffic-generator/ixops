package cluster

import (
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
