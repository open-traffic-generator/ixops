package topology

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg := config.Get()
		// os.Setenv("ENV_IXIA_C_TOPO_TYPE", cfg.IxiaC.KneTopology)
		// ixexec.ExecCmd("dt")
		return nil
	},
}
