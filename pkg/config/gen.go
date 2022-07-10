package config

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var cfg Config

func Gen() error {
	cfg, _ = ReadConfigYaml("/home/ashukuma/conf.yaml")
	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Printf("error while Marshaling. %v\n", err)
		return fmt.Errorf(fmt.Sprintf("error while Marshaling. %v\n", err))
	}

	err = ioutil.WriteFile("/home/ashukuma/config.yaml", yamlData, 0666)
	if err != nil {
		log.Printf("error while wring to %s: %v", "conf.yaml", err)
		return fmt.Errorf(fmt.Sprintf("error while wring to %s: %v", "conf.yaml", err))
	}

	log.Println("Config generated at /home/ashukuma/config.yaml")

	return nil
}

func Get() *Config {
	c, e := ReadConfigYaml("/home/ashukuma/config.yaml")
	if e != nil {
		log.Fatal(e)
	}
	return &c
}
