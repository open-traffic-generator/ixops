package cmd

import (
	"fmt"
	"os"

	"github.com/open-traffic-generator/ixops/cmd/cluster"
	"github.com/open-traffic-generator/ixops/cmd/config"
	"github.com/open-traffic-generator/ixops/cmd/diagnostics"
	"github.com/open-traffic-generator/ixops/cmd/images"
	"github.com/open-traffic-generator/ixops/cmd/otg"
	"github.com/open-traffic-generator/ixops/cmd/tests"
	"github.com/open-traffic-generator/ixops/cmd/topology"
	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ixops",
	Short: "Ixia-C Operations - the easiest way to manage emulated network topologies involving Ixia-C",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		configs.Configure()
		return nil
	},
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	c := configs.GetCmdConfig()
	rootCmd.PersistentFlags().BoolVarP(&c.Debug, "verbose", "v", true, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&c.Quiet, "quiet", "q", false, "Disable logging")
	rootCmd.PersistentFlags().StringVarP(&c.Config, "config", "c", c.Config, "Path to ixops config")

	rootCmd.AddCommand(cluster.Cmd())
	rootCmd.AddCommand(topology.Cmd())
	rootCmd.AddCommand(images.Cmd())
	rootCmd.AddCommand(diagnostics.Cmd())
	rootCmd.AddCommand(tests.Cmd())
	rootCmd.AddCommand(config.Cmd())
	rootCmd.AddCommand(otg.Cmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
