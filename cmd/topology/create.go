package topology

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
