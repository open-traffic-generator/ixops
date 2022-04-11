package cluster

import (
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage Cluster",
}

func Cmd() *cobra.Command {
	clusterCmd.AddCommand(setupCmd)
	clusterCmd.AddCommand(teardownCmd)
	return clusterCmd
}
