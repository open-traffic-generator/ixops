package cfg

import (
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:          "cfg",
	Short:        "Manage cfg",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	cfgCmd.AddCommand(genCmd)
	return cfgCmd
}
