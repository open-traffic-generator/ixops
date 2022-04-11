package cluster

import (
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown Cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
