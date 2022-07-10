package images

import (
	"github.com/open-traffic-generator/ixops/internal/get_images"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get container images",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := get_images.GetImages()
		if err != nil {
			return err
		}
		return nil
	},
}
