package cluster

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/pkg/config"
	"github.com/open-traffic-generator/ixops/pkg/ixexec"
	"github.com/spf13/cobra"
)

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Teardown cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if cfg.ClusterType == "kind" {
			ixexec.ExecCmd("rm_kc")
		} else if cfg.ClusterType == "gcp" {
			ixexec.ExecCmd("rm_gc")
		} else {
			return fmt.Errorf("unsupported cluster type %v", cfg.ClusterType)
		}

		return nil
	},
}
