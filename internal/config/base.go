package config

type IxiaCTestClientRepoConfig struct {
	Url    string `default:"http://gitlab.it.keysight.com/athena/tests.git" yaml:"url"`
	Commit string `default:"821272cd" yaml:"commit"`
}

type IxiaCTestClientConfig struct {
	Yaml string                    `default:"ixia-c-test-client.yaml" yaml:"yaml"`
	Home string                    `default:"tests" yaml:"home"`
	Repo IxiaCTestClientRepoConfig `yaml:"repo"`
}

type IxiaCContainerLabConfig struct {
	Topology string `default:"topo.clab.yaml" yaml:"topology"`
	Tests    string `default:"clabtests" yaml:"tests"`
}

type IxiaCImageConfig struct {
	Controller     string `default:"ghcr.io/open-traffic-generator/licensed/ixia-c-controller" yaml:"controller"`
	TrafficEngine  string `default:"ghcr.io/open-traffic-generator/ixia-c-traffic-engine" yaml:"traffic_engine"`
	ProtocolEngine string `-default:"ghcr.io/open-traffic-generator/licensed/ixia-c-protocol-engine" yaml:"protocol_engine"`
	Grpc           string `default:"ghcr.io/open-traffic-generator/ixia-c-grpc-server" yaml:"grpc"`
	Gnmi           string `default:"ghcr.io/open-traffic-generator/ixia-c-gnmi-server" yaml:"gnmi"`
}

type IxiaCConfig struct {
	Home            string                  `default:"ixia-c" yaml:"home"`
	OperatorVersion string                  `default:"v0.1.94" yaml:"operator_version"`
	Release         string                  `default:"0.0.1-2969" yaml:"release"`
	DutImage        string                  `default:"us-central1-docker.pkg.dev/kt-nts-athena-dev/keysight/ceos:4.28.01F" yaml:"dut_image"`
	KneTopology     string                  `default:"otg-dut-otg" yaml:"kne_topology"`
	TestClient      IxiaCTestClientConfig   `yaml:"test_client"`
	ContainerLab    IxiaCContainerLabConfig `yaml:"containerlab"`
}

type FeaturesProfileConfig struct {
	Repo          string `default:"https://github.com/open-traffic-generator/featureprofiles.git" yaml:"repo"`
	Commit        string `default:"be0a279" yaml:"commit"`
	Home          string `default:"featureprofiles" yaml:"home"`
	KneBindConfig string `default:"kne-bind.yaml" yaml:"kne_bind_config"`
	TestBed       string `default:"topologies/atedut_2.testbed" yaml:"testbed"`
}

type MeshnetConfig struct {
	Repo   string `default:"https://github.com/networkop/meshnet-cni" yaml:"repo"`
	Commit string `default:"de89b2e" yaml:"commit"`
	Home   string `default:"meshnet-cni" yaml:"home"`
	Image  string `default:"networkop/meshnet:v0.3.0" yaml:"image"`
}

type KindConfig struct {
	Version   string `default:"v0.13.0" yaml:"version"`
	NodeCount string `default:"1" yaml:"node_count"`
}

type GCloudNodeConfig struct {
	MasterNodeType  string `default:"e2-standard-4" yaml:"_node_type"`
	WorkerNodeType  string `default:"e2-standard-8" yaml:"worker_node_type"`
	WorkerNodeCount string `default:"1" yaml:"worker_node_count"`
}

type GCloudConfig struct {
	Version       string           `default:"383.0.1" yaml:"version"`
	Home          string           `default:"google-cloud-sdk" yaml:"home"`
	Verbosity     string           `default:"warning" yaml:"verbosity"`
	Account       string           `default:"ixia-c-automation@kt-nts-athena-dev.iam.gserviceaccount.com" yaml:"account"`
	Project       string           `default:"kt-nts-athena-dev" yaml:"project"`
	Region        string           `default:"us-central1" yaml:"region"`
	Zone          string           `default:"us-central1-a" yaml:"zone"`
	ServiceKey    string           `default:"ixia-c-automation.json" yaml:"service_key"`
	EMail         string           `default:"" yaml:"email"`
	ExecutionID   string           `default:"01" yaml:"execution_id"`
	Node          GCloudNodeConfig `yaml:"node"`
	Topology      string           `default:"private" yaml:"topology"`
	Image         string           `default:"ubuntu-pro-2204-jammy-v20220506" yaml:"image"`
	KubeConfigTTL string           `default:"168h0m0s" yaml:"kube_config_ttl"`
}

type KopsConfig struct {
	Version   string `default:"v1.23.1" yaml:"version"`
	Verbosity int    `default:"0" yaml:"verbosity"`
}

type Config struct {
	GoVersion           string                `default:"1.18" yaml:"go_version"`
	ProtoCVersion       string                `default:"3.20.1" yaml:"protoc_version"`
	KubernetsVersion    string                `default:"v1.23.6" yaml:"kubernets_version"`
	MetaLLBVersion      string                `default:"v0.12" yaml:"metallb_version"`
	ContainerLabVersion string                `default:"0.26.2" yaml:"containerlab_version"`
	IxiaC               IxiaCConfig           `yaml:"ixia_c"`
	FeaturesProfile     FeaturesProfileConfig `yaml:"features_profile"`
	Meshnet             MeshnetConfig         `yaml:"meshnet"`
	Kind                KindConfig            `yaml:"kind"`
	Kops                KopsConfig            `yaml:"kops"`
	GCloud              GCloudConfig          `yaml:"gcloud"`
}
