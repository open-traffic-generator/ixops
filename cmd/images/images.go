package images

import (
	"github.com/spf13/cobra"
)

var imagesCmd = &cobra.Command{
	Use:          "images",
	Short:        "Manage container images",
	SilenceUsage: true,
}

func Cmd() *cobra.Command {
	imagesCmd.AddCommand(getCmd)
	imagesCmd.AddCommand(rmCmd)
	return imagesCmd
}
