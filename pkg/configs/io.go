package configs

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func LoadFromBytes(bytes []byte) (*AppConfig, error) {
	c := AppConfig{}

	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return nil, fmt.Errorf("could not unmarshal config bytes: %v", err)
	}

	if err := defaults.Set(&c); err != nil {
		return nil, fmt.Errorf("could not set defaults for config: %v", err)
	}

	return &c, nil
}

func LoadFromFile(path string) (*AppConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", path, err)
	}

	return LoadFromBytes(b)
}

func DumpToFile(c *AppConfig, path string) error {
	b, err := yaml.Marshal(c)

	if err != nil {
		return fmt.Errorf("could not marshal config: %v", err)
	}
	if err := os.WriteFile(path, b, 0666); err != nil {
		return fmt.Errorf("could not write file %s: %v", path, err)
	}

	return nil
}

func GetAppConfig() *AppConfig {
	if appConfig == nil {
		if _, err := os.Stat(cmdConfig.Config); err == nil {
			log.Debug().Str("path", cmdConfig.Config).Msg("Using config")
			appConfig, err = LoadFromFile(cmdConfig.Config)
			if err != nil {
				log.Fatal().Err(err).Msg("Could not load config")
			}
		} else if errors.Is(err, os.ErrNotExist) {
			appConfig, err = LoadFromBytes([]byte{})
			if err != nil {
				log.Fatal().Err(err).Msg("Could not load config from bytes")
			}
			if err := DumpToFile(appConfig, cmdConfig.Config); err != nil {
				log.Fatal().Err(err).Msg("Could not generate config")
			}
			log.Debug().Str("path", cmdConfig.Config).Msg("Generated config")
		} else {
			log.Fatal().Err(err).Str("config", cmdConfig.Config).Msg("Could not stat file")
		}
	}
	return appConfig
}

func GetCmdConfig() *CmdConfig {
	if cmdConfig == nil {
		h, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not determine user home directory")
		}

		home := path.Join(h, ".ixops")

		if err := os.MkdirAll(home, 0777); err != nil {
			log.Fatal().Err(err).Str("home", home).Msg("Could not creat ixops home directory")
		}

		cmdConfig = &CmdConfig{
			Home:   home,
			Debug:  false,
			Config: path.Join(home, "config.yaml"),
		}
	}

	return cmdConfig
}
