package configs

import (
	"gopkg.in/yaml.v2"
)

var cmdConfig *CmdConfig
var appConfig *AppConfig

type CmdConfig struct {
	Home   string `yaml:"home"`
	Debug  bool   `yaml:"debug"`
	Config string `yaml:"config"`
}
type AppConfig struct {
	IxiaC *IxiaC  `yaml:"ixia_c" default:"{}"`
	Nodes []*Node `yaml:"nodes" default:"[{\"master\": true}]"`
}

type IxiaC struct {
	Release string `yaml:"release" default:"latest"`
}

type Node struct {
	Ip     string `yaml:"ip" default:"localhost"`
	Port   uint   `yaml:"port" default:"22"`
	Master bool   `yaml:"master" default:"false"`
}

func (c *CmdConfig) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

func (c *AppConfig) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return err.Error()
	}

	return string(b)
}
