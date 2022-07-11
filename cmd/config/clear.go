package config

import (
	"fmt"
	"os"

	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear config",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := configs.GetCmdConfig()
		if err := os.RemoveAll(c.Config); err != nil {
			return fmt.Errorf("could not clear config")
		}
		return nil
	},
}
