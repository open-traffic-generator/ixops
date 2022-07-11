package config

import (
	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get config",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msgf("Config at %s:\n\n%v", configs.GetCmdConfig().Config, configs.GetAppConfig())
		return nil
	},
}
