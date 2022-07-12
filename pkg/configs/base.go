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
	Topologies *[]*Topology `yaml:"topologies"`
	Endpoints  *[]*Endpoint `yaml:"endpoints"`
	Nodes      *[]*Node     `yaml:"nodes"`
}

func (v *AppConfig) SetDefaults() {
	if v.Topologies == nil {
		p := &Topology{}
		v.Topologies = &[]*Topology{p}
	}
	for _, p := range *v.Topologies {
		p.SetDefaults()
	}
	if v.Endpoints == nil {
		p := &Endpoint{}
		p.SetDefaults()
		v.Endpoints = &[]*Endpoint{p}
	}
	for _, p := range *v.Endpoints {
		p.SetDefaults()
	}
	if v.Nodes == nil {
		p := &Node{}
		p.SetDefaults()
		v.Nodes = &[]*Node{p}
	}
	for _, p := range *v.Nodes {
		p.SetDefaults()
	}
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
