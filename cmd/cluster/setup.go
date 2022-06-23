package cluster

import (
	"github.com/open-traffic-generator/ixops/internal/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := setup.CommonSetup(&args)
		if err != nil {
			return err
		}
		err = setup.SetupCluster()
		if err != nil {
			return err
		}
		return nil
	},
}
