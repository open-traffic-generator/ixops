package topology

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		// cfg := config.Get()
		// os.Setenv("ENV_IXIA_C_TOPO_TYPE", cfg.IxiaC.KneTopology)
		// ixexec.ExecCmd("ct")
		return nil
	},
}
