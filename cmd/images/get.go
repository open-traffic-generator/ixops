package images

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get container images",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
