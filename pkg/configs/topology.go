package configs

type TopologyPlatform string

const (
	TopologyPlatformDocker TopologyPlatform = "docker"
	TopologyPlatformKind   TopologyPlatform = "kind"
)

func (t *TopologyPlatform) IsValid() bool {
	switch *t {
	case TopologyPlatformDocker:
		return true
	case TopologyPlatformKind:
		return true
	}

	return false
}

type Topology struct {
	Name        *string           `yaml:"name"`
	CreateLinks *bool             `yaml:"create_links"`
	Platform    *TopologyPlatform `yaml:"platform"`
	Links       *[]string         `yaml:"links"`
}

func (v *Topology) SetDefaults() {
	SetDefaultString(&v.Name, "otg-b2b")
	SetDefaultBool(&v.CreateLinks, true)
	if v.Platform == nil {
		p := TopologyPlatform("docker")
		v.Platform = &p
	}
	if v.Links == nil {
		v.Links = &[]string{"otg:veth1 otg:veth2"}
	}
}
