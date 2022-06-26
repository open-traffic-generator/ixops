package tests

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/internal/config"
	"github.com/open-traffic-generator/ixops/internal/get_tests"
	"github.com/open-traffic-generator/ixops/internal/utils"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("please pass config file as argument \n Ex:ixops tests get conf.yaml")
		} else {
			if utils.FileExists(args[0]) {
				conf, err := config.ReadConfigYaml(args[0])
				if err != nil {
					return err
				}
				get_tests.GetTests(conf)
				return nil
			} else {
				return fmt.Errorf("config file doesn't exists")
			}
		}
	},
}
