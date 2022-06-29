package cluster

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/internal/config"
	"github.com/open-traffic-generator/ixops/internal/setup"
	"github.com/open-traffic-generator/ixops/internal/utils"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args[0]) != 0 {
			kopsSetup := true
			if len(args) > 1 && args[1] == "kind" {
				kopsSetup = false
			}
			if utils.FileExists(args[0]) {
				_, err := config.ReadConfigYaml(args[0])
				if err != nil {
					return err
				}
				err = setup.CommonSetup(&args)
				if err != nil {
					return err
				}
				err = setup.SetupCluster(kopsSetup)
				if err != nil {
					return err
				}

				return nil
			} else {
				return fmt.Errorf("config file doesn't exists")
			}
		} else {
			return fmt.Errorf("config file should be provided")
		}
	},
}
