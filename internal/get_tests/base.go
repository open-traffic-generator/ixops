package get_tests

type TcPodMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type TcPodSpec struct {
	Containers []Container `yaml:"containers"`
}

type Container struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Command []string `yaml:"command"`
}

type TcPodConfig struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   TcPodMetadata `yaml:"metadata"`
	Spec       TcPodSpec     `yaml:"spec"`
}
