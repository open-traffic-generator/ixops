package topology

import (
	"github.com/spf13/cobra"
)

var topologyCmd = &cobra.Command{
	Use:          "topology",
	Short:        "Manage topology",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	topologyCmd.AddCommand(createCmd)
	topologyCmd.AddCommand(deleteCmd)
	return topologyCmd
}
