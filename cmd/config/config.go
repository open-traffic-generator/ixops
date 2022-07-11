package config

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:          "config",
	Short:        "Manage config",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	configCmd.AddCommand(getCmd)
	configCmd.AddCommand(clearCmd)
	return configCmd
}
