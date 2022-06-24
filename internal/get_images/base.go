package get_images

const (
	Ixiacversion  = "0.0.1-2969"
	FileName      = "ixia-configmap.yaml"
	centralDocker = "us-central1-docker.pkg.dev/kt-nts-athena-dev/keysight/"
	ghrc          = "ghcr.io/open-traffic-generator/"
	Pat           = "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"
	OtgGit        = "https://github.com/open-traffic-generator/ixia-c/releases/download"
)

var configMapVersions string

type ConfigMap struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   ConfigMapMetadata `yaml:"metadata"`
	Data       ConfigMapData     `yaml:"data"`
}

type ConfigMapMetadata struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
}

type ConfigMapData struct {
	Versions string `yaml:"versions"`
}

type Version struct {
	Release string  `json:"release"`
	Images  []Image `json:"images"`
}

type Image struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Tag  string `json:"tag"`
}
