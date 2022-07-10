package cfg

import (
	"github.com/open-traffic-generator/ixops/pkg/config"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Get cfg",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.Gen()
		return err
	},
}
