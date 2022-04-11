package images

import (
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove container images",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
