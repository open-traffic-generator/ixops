package otg

import (
	"github.com/open-traffic-generator/ixops/pkg/otg"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate traffic",
	RunE: func(cmd *cobra.Command, args []string) error {
		oc := otg.NewOtgConfig()
		otg.Generate(oc)
		return nil
	},
}
