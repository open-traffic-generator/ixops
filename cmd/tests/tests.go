package tests

import (
	"github.com/spf13/cobra"
)

var testsCmd = &cobra.Command{
	Use:          "tests",
	Short:        "Get and Run tests",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	testsCmd.AddCommand(getCmd)
	testsCmd.AddCommand(runCmd)
	testsCmd.AddCommand(rmCmd)
	return testsCmd
}
