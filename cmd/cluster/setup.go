package cluster

import (
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg := config.Get()
		// if cfg.ClusterType == "kind" {
		// 	ixexec.ExecCmd("new_kc")
		// } else if cfg.ClusterType == "gcp" {
		// 	ixexec.ExecCmd("new_gc")
		// } else {
		// 	return fmt.Errorf("unsupported cluster type %v", cfg.ClusterType)
		// }

		return nil
	},
}
