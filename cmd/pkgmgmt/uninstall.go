package pkgmgmt

import (
	"github.com/open-traffic-generator/ixops/pkg/pkgmgmt"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [package name]",
	Short: "Uninstall Packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkgmgmt.UninstallPackage(args[0])
	},
}
