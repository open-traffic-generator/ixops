package configs

type IxiaC struct {
	Release       *string `yaml:"release"`
	Free          *bool   `yaml:"free"`
	Controller    *string `yaml:"controller"`
	TrafficEngine *string `yaml:"traffic_engine"`
}

func (v *IxiaC) SetDefaults() {
	SetDefaultString(&v.Release, "0.0.1-2994")
	SetDefaultBool(&v.Free, true)
	SetDefaultString(&v.Controller, "ghcr.io/open-traffic-generator/ixia-c-controller:0.0.1-2994")
	SetDefaultString(&v.TrafficEngine, "ghcr.io/open-traffic-generator/ixia-c-traffic-engine:1.4.1.29")
}
