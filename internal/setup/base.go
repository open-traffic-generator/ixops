package setup

type KindNetworkInfo struct {
	APIServerAddress string `yaml:"apiServerAddress"`
	APIServerPort    int    `yaml:"apiServerPort"`
}

type KindRoleInfo struct {
	Role string `yaml:"role"`
}

type KindConfig struct {
	Kind       string          `yaml:"kind"`
	APIVersion string          `yaml:"apiVersion"`
	Networking KindNetworkInfo `yaml:"networking"`
	Nodes      []KindRoleInfo  `yaml:"nodes"`
}

type MetallbMetadata struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
}

type MetallbData struct {
	Config string `yaml:"config"`
}

type MetallbConfig struct {
	APIVersion string          `yaml:"apiVersion"`
	Kind       string          `yaml:"kind"`
	Metadata   MetallbMetadata `yaml:"metadata"`
	Data       MetallbData     `yaml:"data"`
}

const (
	KindConfigFile       = "kind.yaml"
	KindVersion          = "v0.13.0"
	NodeCount            = 1
	TimeOut              = 300
	MetallbVersion       = "v0.12"
	MetallbConfigFile    = "metallb.yaml"
	IxiaCOperatorVersion = "v0.1.94"
	MeshnetCommit        = "de89b2e"
	MeshnetVersion       = "v0.3.0"
)
