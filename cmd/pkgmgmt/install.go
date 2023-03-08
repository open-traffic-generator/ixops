package pkgmgmt

import (
	"github.com/open-traffic-generator/ixops/pkg/pkgmgmt"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [package name]",
	Short: "Install Packages",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkgmgmt.InstallPackage(args[0], "")
	},
}
