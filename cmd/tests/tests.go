package tests

import (
	"github.com/spf13/cobra"
)

var testsCmd = &cobra.Command{
	Use:          "tests",
	Short:        "Manage tests",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	testsCmd.AddCommand(getCmd)
	testsCmd.AddCommand(rmCmd)
	testsCmd.AddCommand(runCmd)
	testsCmd.AddCommand(lsCmd)
	return testsCmd
}
