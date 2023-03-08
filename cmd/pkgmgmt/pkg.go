package pkgmgmt

import (
	"github.com/spf13/cobra"
)

var pkgCmd = &cobra.Command{
	Use:          "pkg",
	Short:        "Manage Packages",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	pkgCmd.AddCommand(installCmd)
	pkgCmd.AddCommand(uninstallCmd)
	return pkgCmd
}
