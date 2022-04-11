package topology

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
