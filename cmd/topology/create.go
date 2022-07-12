package topology

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/open-traffic-generator/ixops/pkg/dockerc"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := configs.GetAppConfig()
		if len(*c.Topologies) != 1 {
			return fmt.Errorf("exactly one topology needs to be specified")
		}

		t := (*c.Topologies)[0]
		switch p := t.Platform; *p {
		case configs.TopologyPlatformDocker:
			return dockerc.CreateTopology(t)
		default:
			return fmt.Errorf("topology platform %s not supported", *p)
		}
	},
}
