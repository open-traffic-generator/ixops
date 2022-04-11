package diagnostics

import (
	"github.com/spf13/cobra"
)

var diagnosticsCmd = &cobra.Command{
	Use:          "diagnostics",
	Short:        "Manage diagnostics",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	diagnosticsCmd.AddCommand(collectCmd)
	return diagnosticsCmd
}
