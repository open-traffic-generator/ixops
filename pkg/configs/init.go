package configs

import "github.com/rs/zerolog"

func init() {
	GetCmdConfig()
	InitLogger()
}

func Configure() {
	if cmdConfig.Debug {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if cmdConfig.Quiet {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
