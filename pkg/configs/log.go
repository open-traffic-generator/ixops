package configs

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}
