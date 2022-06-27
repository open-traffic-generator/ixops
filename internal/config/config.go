package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

func ReadConfigYaml(filePath string) (Config, error) {
	config := Config{}
	errorString := ""
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		errorString = fmt.Sprintf("failed to read %s: %v", filePath, err)
		log.Println(err)
		return config, fmt.Errorf(errorString)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		errorString = fmt.Sprintf("failed to unmarshall %v: %v", string(yamlFile), err)
		log.Println(err)
		return config, fmt.Errorf(errorString)
	}
	return config, nil
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(c)

	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}
