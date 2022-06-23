package topology

import (
	"strings"

	"github.com/open-traffic-generator/ixops/internal/topology"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 || args == nil {
			args = []string{"otg-dut-otg"}
		}
		if strings.Contains(args[0], ".yaml") || strings.Contains(args[0], ".txt") {
			err := topology.CreateTopologyWithFile(args[0])
			if err != nil {
				return err
			}
		} else {
			err := topology.CreateTopologyWithTopoType(topology.TopologyType(args[0]))
			if err != nil {
				return err
			}
		}
		return nil
	},
}
