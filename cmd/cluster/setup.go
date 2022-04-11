package cluster

import (
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
