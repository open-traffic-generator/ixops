package configs

import (
	"gopkg.in/yaml.v2"
)

var cmdConfig *CmdConfig
var appConfig *AppConfig

type CmdConfig struct {
	Home   string `yaml:"home"`
	Debug  bool   `yaml:"debug"`
	Quiet  bool   `yaml:"quiet"`
	Config string `yaml:"config"`
}
type AppConfig struct {
	IxiaC *IxiaC  `yaml:"ixia_c" default:"{}"`
	Nodes []*Node `yaml:"nodes" default:"[{\"master\": true}]"`
}

type IxiaC struct {
	Release string `yaml:"release" default:"latest"`
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
