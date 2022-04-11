package topology

import (
	"github.com/spf13/cobra"
)

var topologyCmd = &cobra.Command{
	Use:   "topology",
	Short: "Manage Topology",
}

func Cmd() *cobra.Command {
	topologyCmd.AddCommand(createCmd)
	topologyCmd.AddCommand(deleteCmd)
	return topologyCmd
}
