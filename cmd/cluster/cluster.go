package cluster

import (
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:          "cluster",
	Short:        "Manage cluster",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	clusterCmd.AddCommand(setupCmd)
	clusterCmd.AddCommand(teardownCmd)
	return clusterCmd
}
