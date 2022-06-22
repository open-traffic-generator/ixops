package cluster

import (
	"github.com/open-traffic-generator/ixops/internal/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup.CreateKindConfig()
		// setup.InstallKind("0.13.0")
		// setup.CreateKindCluster()
		// setup.GetKubectl()
		// setup.GetMetallb()
		setup.MakeMetallbConfig()
		return nil
	},
}
