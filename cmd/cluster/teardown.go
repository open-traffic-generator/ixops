package cluster

import (
	"github.com/open-traffic-generator/ixops/internal/setup"
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		gcCluster := true
		if len(args) > 0 && args[0] == "kind" {
			gcCluster = false
		}

		return setup.TeardownCluster(gcCluster)
	},
}
