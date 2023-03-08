package pkgmgmt

import (
	"github.com/open-traffic-generator/ixops/pkg/pkgmgmt"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [package name] [version]",
	Short: "Install Packages",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := ""
		if len(args) == 2 {
			version = args[1]
		}
		return pkgmgmt.InstallPackage(args[0], version)
	},
}
