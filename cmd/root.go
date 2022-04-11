package cmd

import (
	"fmt"
	"os"

	"github.com/open-traffic-generator/ixops/cmd/cluster"
	"github.com/open-traffic-generator/ixops/cmd/images"
	"github.com/open-traffic-generator/ixops/cmd/topology"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

var rootCmd = &cobra.Command{
	Use:          "ixops",
	Short:        "Ixia-C Operations - the easiest way to manage emulated network topologies involving Ixia-C",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(cluster.Cmd())
	rootCmd.AddCommand(topology.Cmd())
	rootCmd.AddCommand(images.Cmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "Enable verbose logging")
}
