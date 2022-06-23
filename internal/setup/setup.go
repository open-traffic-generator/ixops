package setup

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
	"gopkg.in/yaml.v2"
)

func createKindConfig(configFile string, nodes int) (string, error) {
	kindConfig := KindConfig{
		Kind:       "Cluster",
		APIVersion: "kind.x-k8s.io/v1alpha4",
		Networking: KindNetworkInfo{
			APIServerAddress: "127.0.0.1",
			APIServerPort:    6443,
		},
		Nodes: []KindRoleInfo{},
	}
	kindConfig.Nodes = append(kindConfig.Nodes, KindRoleInfo{Role: "control-plane"})
	if nodes > 2 {
		log.Printf("Node Count: %d", nodes)
		kindConfig.Nodes = append(kindConfig.Nodes, KindRoleInfo{Role: "worker"})
	}
	yamlData, err := yaml.Marshal(&kindConfig)
	if err != nil {
		log.Printf("error while Marshaling. %v\n", err)
		return "", fmt.Errorf(fmt.Sprintf("error while Marshaling. %v\n", err))
	}

	commonHome, err := utils.GetCommonHome()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(commonHome, configFile)
	err = ioutil.WriteFile(filePath, yamlData, 0666)
	if err != nil {
		log.Printf("error while wring to %s: %v", configFile, err)
		return "", fmt.Errorf(fmt.Sprintf("error while wring to %s: %v", configFile, err))
	}
	return filePath, nil
}

func getKind(kindVersion string) error {
	log.Printf("Installing kind@%s", kindVersion)
	_, err := utils.ExecCmd("go", "install", fmt.Sprintf("sigs.k8s.io/kind@%v", kindVersion))
	if err != nil {
		log.Printf("failed to get kind: %v", err)
		return fmt.Errorf(fmt.Sprintf("failed to get kind: %v", err))
	}
	return nil
}

func kindClusterExists() (bool, error) {
	errorString := ""
	log.Println("Checking for existing kind cluster....")
	_, err := utils.ExecCmd("kind", "get", "clusters")
	if err != nil {
		if strings.Contains(err.Error(), "No kind clusters found") {
			log.Println("No existing kind cluster...")
			return false, nil
		} else {
			errorString = fmt.Sprintf("Failed checking for any existing kind cluster: %v", err)
			log.Println(errorString)
			return false, err
		}
	}
	log.Println("kind cluster already exists...")
	return true, nil
}

func deployBasicCluster(configFilePath string, waitTime int) error {
	log.Printf("Deploying basic kind cluster...\n")
	_, err := utils.ExecCmd("kind", "create", "cluster", "--config", configFilePath, "--wait", fmt.Sprintf("%ds", waitTime))
	if err != nil {
		log.Printf("failed to deploy basic kind cluster: %v", err)
		return fmt.Errorf(fmt.Sprintf("failed to deploy basic kind cluster: %v", err))
	}
	return nil
}

func getKubectl() error {
	errorString := ""

	log.Printf("Removing older kubectl...")
	_, err := utils.ExecCmd("rm", "-rf", "kubectl")
	if err != nil {
		errorString = fmt.Sprintf("Failed to remove older kubectl: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Printf("Copying kubectl from kind container....")
	_, err = utils.ExecCmd("docker", "cp", "kind-control-plane:/usr/bin/kubectl", "./")
	if err != nil {
		errorString = fmt.Sprintf("Failed to copy kubectl from kind container: %v", err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	log.Printf("Installing kubectl....")
	_, err = utils.ExecCmd("sudo", "install", "-o", "root", "-g", "root", "-m", "0755", "kubectl", "/usr/local/bin/kubectl")
	if err != nil {
		errorString = fmt.Sprintf("Failed to install kubectl: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Printf("Removing copied kubectl...")
	_, err = utils.ExecCmd("rm", "-rf", "kubectl")
	if err != nil {
		errorString = fmt.Sprintf("Failed to remove copied kubectl: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	return nil
}

func createMetallbConfig(configFile string) (string, error) {
	errorString := ""
	log.Println("Getting docker network information...")
	out, err := utils.ExecCmd("docker", "network", "inspect", "-f", "{{.IPAM.Config}}", "kind")
	if err != nil {
		errorString = fmt.Sprintf("Failed to get docker network information: %v", err)
		log.Println(err)
		return "", fmt.Errorf(errorString)
	}
	address := strings.Split(strings.Split(strings.Split(out, " ")[0], "{")[1], "/")[0]
	prefix := strings.Join(strings.Split(address, ".")[:3], ".")
	log.Printf("Prefix: %s", prefix)

	log.Println("Creating metallb config...")
	metallbConfig := MetallbConfig{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: MetallbMetadata{
			Namespace: "metallb-system",
			Name:      "config",
		},
		Data: MetallbData{
			Config: "",
		},
	}

	addressInfo := fmt.Sprintf("	  - %s.100 - %s.250", prefix, prefix)
	metallbConfig.Data.Config = "address-pools:\n" +
		" - name: default\n" +
		"	protocol: layer2\n" +
		"	addresses:\n" +
		addressInfo

	yamlData, err := yaml.Marshal(&metallbConfig)
	if err != nil {
		errorString = fmt.Sprintf("Error while Marshaling. %v", err)
		log.Println(err)
		return "", fmt.Errorf(errorString)
	}
	log.Printf("Metallb Config: %s", string(yamlData))

	commonHome, err := utils.GetCommonHome()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(commonHome, configFile)
	err = ioutil.WriteFile(filePath, yamlData, 0666)
	if err != nil {
		errorString = fmt.Sprintf("Error while wring to %s: %v", filePath, err)
		log.Println(err)
		return "", fmt.Errorf(errorString)
	}

	return filePath, nil

}

func getMetallb(version string, metallbConfigFile string, waitTime int64) error {
	errorString := ""
	log.Printf("Apply metallb namespace.yaml...")
	_, err := utils.ExecCmd("kubectl", "apply", "-f", fmt.Sprintf("https://raw.githubusercontent.com/metallb/metallb/%s/manifests/namespace.yaml", version))
	if err != nil {
		errorString = fmt.Sprintf("failed to apply metallb namespace.yaml: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Printf("Creating secrets for metallb...")
	_, err = utils.ExecCmd("kubectl", "create", "generic", "-n", "metallb-system", "memberlist", "--from-literal=secretkey=\"$(openssl rand -base64 128)\"")
	if err != nil {
		errorString = fmt.Sprintf("failed to creating secrets for metallbs: %v", err)
		log.Println(errorString)
	}

	log.Printf("Apply metallb.yaml...")
	utils.ExecCmd("kubectl", "apply", "-f", fmt.Sprintf("https://raw.githubusercontent.com/metallb/metallb/%s/manifests/metallb.yaml", version))

	log.Printf("Waiting for pods to be ready...")
	kubeClient, err := utils.NewK8sClient()
	if err != nil {
		errorString = fmt.Sprintf("failed to create k8s client: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	err = utils.WaitFor(func() (bool, error) { return kubeClient.AllPodsAreReady("metallb-system") }, nil)
	if err != nil {
		errorString = fmt.Sprintf("metallb-system pods are not ready: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	metallbConfigFilePath, err := createMetallbConfig(metallbConfigFile)
	if err != nil {
		return err
	}

	log.Printf("Applying metallb config...")
	_, err = utils.ExecCmd("kubectl", "apply", "-f", metallbConfigFilePath)
	if err != nil {
		errorString = fmt.Sprintf("failed to apply metallb config: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	return nil
}

func getMeshnet(commit string, version string, waitTime int64) error {
	errorString := ""

	log.Println("removing older mershnet-cni...")
	_, err := utils.ExecCmd("rm", "-rf", "meshnet-cni")
	if err != nil {
		errorString = fmt.Sprintf("failed to remove older meshnet-cni: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Println("Cloning mershnet repo...")
	utils.ExecCmd("git", "clone", "https://github.com/networkop/meshnet-cni")

	currentDirectory, err := os.Getwd()
	if err != nil {
		errorString = fmt.Sprintf("failed get curent working directory: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	log.Printf("Current working directory : %v", currentDirectory)

	meshnetDirectory := filepath.Join(currentDirectory, "meshnet-cni")
	log.Printf("changing working directory to: %s", meshnetDirectory)
	err = os.Chdir(meshnetDirectory)
	if err != nil {
		errorString = fmt.Sprintf("failed change working directory to %s: %v", meshnetDirectory, err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Printf("checking out to meshnet-cni commit %s\n", commit)
	utils.ExecCmd("git", "checkout", commit)

	log.Printf("changing working directory to: %s", currentDirectory)
	err = os.Chdir(currentDirectory)
	if err != nil {
		errorString = fmt.Sprintf("failed change working directory to %s: %v", currentDirectory, err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	dsetFile := filepath.Join(meshnetDirectory, "manifests", "base", "daemonset.yaml")
	log.Printf("updating %s", dsetFile)
	dsetFileInput, err := ioutil.ReadFile(dsetFile)
	if err != nil {
		errorString = fmt.Sprintf("failed to read %s: %v", dsetFile, err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	lines := strings.Split(string(dsetFileInput), "\n")
	for i, line := range lines {
		if strings.Contains(line, "image: networkop/meshnet:latest") {
			lines[i] = "          image: networkop/meshnet:" + version
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(dsetFile, []byte(output), 0644)
	if err != nil {
		errorString = fmt.Sprintf("failed to write %s: %v", dsetFile, err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Println("Applying mesthnet yaml...")
	baseDirectory := filepath.Join(currentDirectory, "meshnet-cni", "manifests", "base")
	utils.ExecCmd("kubectl", "apply", "-k", baseDirectory)

	log.Printf("Waiting for pods to be ready...")
	kubeClient, err := utils.NewK8sClient()
	if err != nil {
		errorString = fmt.Sprintf("failed to create k8s client: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	err = utils.WaitFor(func() (bool, error) { return kubeClient.AllPodsAreReady("meshnet") }, nil)
	if err != nil {
		errorString = fmt.Sprintf("meshnet pods are not ready: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	return nil
}

func getIxiaCOperator(version string, waitTime int64) error {
	errorString := ""

	log.Println("Creating ixia-c-operator....")
	_, err := utils.ExecCmd("kubectl", "apply", "-f", fmt.Sprintf("https://github.com/open-traffic-generator/ixia-c-operator/releases/download/%s/ixiatg-operator.yaml", version))
	if err != nil {
		errorString = fmt.Sprintf("failed to create ixia-c-operator: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	log.Printf("Waiting for pods to be ready...")
	kubeClient, err := utils.NewK8sClient()
	if err != nil {
		errorString = fmt.Sprintf("failed to create k8s client: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	err = utils.WaitFor(func() (bool, error) { return kubeClient.AllPodsAreReady("ixiatg-op-system") }, nil)
	if err != nil {
		errorString = fmt.Sprintf("ixiatg-op-system pods are not ready: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	return nil
}

func SetupCluster() error {
	err := getKind(KindVersion)
	if err != nil {
		return err
	}

	clusterExists, err := kindClusterExists()
	if err != nil {
		return err
	}

	if !clusterExists {
		kindConfigFilePath, err := createKindConfig(KindConfigFile, NodeCount)
		if err != nil {
			return err
		}
		deployBasicCluster(kindConfigFilePath, TimeOut)
	}

	err = getKubectl()
	if err != nil {
		return err
	}

	err = getMetallb(MetallbVersion, MetallbConfigFile, TimeOut)
	if err != nil {
		return err
	}

	err = getMeshnet(MeshnetCommit, MeshnetVersion, TimeOut)
	if err != nil {
		return err
	}

	err = getIxiaCOperator(IxiaCOperatorVersion, TimeOut)
	if err != nil {
		return err
	}
	return nil
}
