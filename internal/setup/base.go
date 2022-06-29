package setup

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
)

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
	KopsVersion          = "v1.23.1"
	KubernetesVersion    = "v1.23.6"
	GCloudVersion        = "383.0.1"
	GCloudVerbosity      = "warning"
	GCloudAccount        = "ixia-c-automation@kt-nts-athena-dev.iam.gserviceaccount.com"
	GCloudProject        = "kt-nts-athena-dev"
	GCloudRegion         = "us-central1"
	GCloudZone           = "us-central1-a"
	GCloudUser           = "test"
	GCloudEmail          = ""
	GCloudWorkerNodes    = 1
	GCloudMasterNodeType = "e2-standard-4"
	GCloudWorkerNodeType = "e2-standard-8"
	GCloudTopology       = "private"
	GCloudNetworking     = "calico"
	GCloudKubeconfigTTL  = "168h0m0s"
	KopsVerbosity        = 0
	NodeCount            = 1
	TimeOut              = 300
	MetallbVersion       = "v0.12"
	MetallbConfigFile    = "metallb.yaml"
	IxiaCOperatorVersion = "v0.1.94"
	MeshnetCommit        = "de89b2e"
	MeshnetVersion       = "v0.3.0"
)

var (
	ClusterTypeGC = true
)

func setKopsEnv() error {
	if err := os.Setenv("KOPS_FEATURE_FLAGS", "AlphaAllowGCE"); err != nil {
		return err
	}
	if err := os.Setenv("CLOUDSDK_COMPUTE_REGION", GCloudRegion); err != nil {
		return err
	}
	if err := os.Setenv("CLOUDSDK_COMPUTE_ZONE", GCloudZone); err != nil {
		return err
	}

	home := os.Getenv("HOME")
	gcSvcAccKey := fmt.Sprintf("%s/.ixops/ixia-c-automation.json", home)
	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", gcSvcAccKey); err != nil {
		return err
	}
	return nil
}

func kopsClusterExists(gcClusterName string, gcStoreName string) (bool, error) {
	if err := setKopsEnv(); err != nil {
		log.Printf("Failed to set environment variables - %v\n", err)
		return false, err
	}

	log.Printf("Checking if cluster %s (store %s) exists\n", gcClusterName, gcStoreName)
	out, err := utils.ExecCmd("kops", "get", "clusters", fmt.Sprintf("--name=%s", gcClusterName), fmt.Sprintf("--state=%s", gcStoreName))
	if err != nil {
		//log.Printf("Failed to check cluster error - %v\n", err)
		return !strings.Contains(err.Error(), "cluster not found"), nil
	}
	return !strings.Contains(out, "cluster not found"), nil
}
