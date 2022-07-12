package images

import (
	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/open-traffic-generator/ixops/pkg/dockerc"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get container images",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := configs.GetAppConfig()
		d, err := dockerc.NewClient((*c.Nodes)[0])
		if err != nil {
			return err
		}
		if err := d.ListImages(); err != nil {
			return err
		}
		return nil
	},
}
