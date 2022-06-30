package topology

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/kne/proto/topo"
	topo1 "github.com/google/kne/topo"
	"github.com/open-traffic-generator/ixops/internal/setup"
	"github.com/open-traffic-generator/ixops/internal/utils"
	"gopkg.in/yaml.v2"
)

func createTopologyFile(topoType TopologyType, topologyFile string) error {
	topology := topo.Topology{
		Name: "ixia-c",
		Nodes: []*topo.Node{
			{
				Name:    "otg",
				Type:    topo.Node_IXIA_TG,
				Version: "0.0.1-2969",
				Services: map[uint32]*topo.Service{
					443: {
						Name:   "https",
						Inside: 443,
					},
					40051: {
						Name:   "grpc",
						Inside: 40051,
					},
					50051: {
						Name:   "gnmi",
						Inside: 50051,
					},
				},
			},
			{
				Name: "dut1",
				Type: topo.Node_ARISTA_CEOS,
				Config: &topo.Config{
					ConfigPath: "/mnt/flash",
					ConfigFile: "startup-config",
					Image:      "ghcr.io/open-traffic-generator/ceos:4.28.01F",
					Cert: &topo.CertificateCfg{
						Config: &topo.CertificateCfg_SelfSigned{
							SelfSigned: &topo.SelfSignedCertCfg{
								CertName: "gnmiCert.pem",
								KeyName:  "gnmiCertKey.pem",
								KeySize:  4096,
							},
						},
					},
					ConfigData: &topo.Config_File{
						File: "./init_dut.txt",
					},
				},
				Services: map[uint32]*topo.Service{
					22: {
						Name:   "ssh",
						Inside: 22,
					},
					6030: {
						Name:   "gnmi",
						Inside: 6030,
					},
				},
			},
		},
		Links: []*topo.Link{},
	}

	if topoType == OTG_DUT_OTG {
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "otg",
			AInt:  "eth1",
			ZNode: "dut1",
			ZInt:  "eth1",
		})
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "otg",
			AInt:  "eth2",
			ZNode: "dut1",
			ZInt:  "eth2",
		})
	} else if topoType == OTG_DUT_DUT_2_OTG {
		topology.Nodes = append(topology.Nodes, &topo.Node{
			Name: "dut2",
			Type: topo.Node_ARISTA_CEOS,
			Config: &topo.Config{
				ConfigPath: "/mnt/flash",
				ConfigFile: "startup-config",
				Image:      "ghcr.io/open-traffic-generator/ceos:4.28.01F",
				Cert: &topo.CertificateCfg{
					Config: &topo.CertificateCfg_SelfSigned{
						SelfSigned: &topo.SelfSignedCertCfg{
							CertName: "gnmiCert.pem",
							KeyName:  "gnmiCertKey.pem",
							KeySize:  4096,
						},
					},
				},
				ConfigData: &topo.Config_File{
					File: "./init_dut.txt",
				},
			},
			Services: map[uint32]*topo.Service{
				22: {
					Name:   "ssh",
					Inside: 22,
				},
				6030: {
					Name:   "gnmi",
					Inside: 6030,
				},
			},
		})
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "otg",
			AInt:  "eth1",
			ZNode: "dut1",
			ZInt:  "eth1",
		})
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "dut1",
			AInt:  "eth2",
			ZNode: "dut2",
			ZInt:  "eth1",
		})
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "dut2",
			AInt:  "eth2",
			ZNode: "otg",
			ZInt:  "eth2",
		})
		topology.Links = append(topology.Links, &topo.Link{
			ANode: "dut2",
			AInt:  "eth2",
			ZNode: "otg",
			ZInt:  "eth3",
		})
	} else {
		log.Printf("unsupported topology type: %v\n", topoType)
		return fmt.Errorf(fmt.Sprintf("unsupported topology type: %v\n", topoType))
	}

	yamlData, err := yaml.Marshal(&topology)
	if err != nil {
		log.Printf("error while Marshaling. %v\n", err)
		return fmt.Errorf(fmt.Sprintf("error while Marshaling. %v\n", err))
	}

	err = ioutil.WriteFile(topologyFile, yamlData, 0666)
	if err != nil {
		log.Printf("error while wring to %s: %v", topologyFile, err)
		return fmt.Errorf(fmt.Sprintf("error while wring to %s: %v", topologyFile, err))
	}
	return nil
}

func createInitDUTConfig(initFile string) error {
	initConfigText := "transceiver qsfp default-mode 4x10G\n" +
		"!\n" +
		"service routing protocols model ribd\n" +
		"!\n" +
		"agent Bfd shutdown\n" +
		"agent PowerManager shutdown\n" +
		"agent LedPolicy shutdown\n" +
		"agent Thermostat shutdown\n" +
		"agent PowerFuse shutdown\n" +
		"agent StandbyCpld shutdown\n" +
		"agent LicenseManager shutdown\n" +
		"!\n" +
		"spanning-tree mode mstp\n" +
		"!\n" +
		"no aaa root\n" +
		"aaa authentication policy local allow-nopassword-remote-login\n" +
		"!\n" +
		"username admin privilege 15 role network-admin nopassword\n" +
		"!\n" +
		"management api gnmi\n" +
		"transport grpc default\n" +
		"    ssl profile octa-ssl-profile\n" +
		"provider eos-native\n" +
		"!\n" +
		"management security\n" +
		"ssl profile octa-ssl-profile\n" +
		"    certificate gnmiCert.pem key gnmiCertKey.pem\n" +
		"!\n" +
		"ip routing\n" +
		"!\n" +
		"ipv6 unicast-routing\n" +
		"!\n" +
		"end"

	err := ioutil.WriteFile(initFile, []byte(initConfigText), 0666)
	if err != nil {
		log.Printf("error while wring to %s: %v", initFile, err)
		return fmt.Errorf(fmt.Sprintf("error while wring to %s: %v", initFile, err))
	}
	return nil
}

func CreateTopologyWithTopoType(topoType TopologyType) error {
	errorString := ""
	err := createTopologyFile(topoType, TopologyFile)
	if err != nil {
		errorString = fmt.Sprintf("Failed create topology File %s: %v", TopologyFile, err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	err = createInitDUTConfig(AristaInitConfigFile)
	if err != nil {
		errorString = fmt.Sprintf("Failed create dut init config file %s: %v", AristaInitConfigFile, err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	err = CreateTopologyWithFile(TopologyFile)
	if err != nil {
		return err
	}
	return nil
}

func getTopologyNamespace(topologyFile string) (string, error) {
	errorString := ""
	data, err := ioutil.ReadFile(topologyFile)
	if err != nil {
		errorString = fmt.Sprintf("failed reading data from file: %s", err)
		log.Println(errorString)
		return "", fmt.Errorf(errorString)
	}
	lines := strings.Split(string(data), "\n")
	if strings.Contains(lines[0], "name:") {
		return strings.TrimSpace(strings.Split(lines[0], "name:")[0]), nil
	}
	errorString = fmt.Sprintf("failed reading namespace from file: %s", err)
	log.Println(errorString)
	return "", fmt.Errorf(errorString)
}

func createSecrets(topoFilePath string) error {
	file, err := os.Open(topoFilePath)
	if err != nil {
		log.Printf("Failed to read topology file %s\n", topoFilePath)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	namespace := ""
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "name:") {
			namespace = strings.TrimSpace(text[5:])
			if namespace[0] == '"' {
				namespace = namespace[1 : len(namespace)-1]
			}
			break
		}
	}
	if namespace == "" {
		errStr := fmt.Sprintf("Failed to find namespace in topology file\n")
		log.Println(errStr)
		return fmt.Errorf(errStr)
	}

	log.Printf("Creating namespace %s\n", namespace)
	_, err = utils.ExecCmd("kubectl", "create", "ns", namespace)
	if err != nil && !strings.Contains(err.Error(), "AlreadyExists") {
		log.Printf("Failed to create namespace error - %v\n", err)
		return err
	}

	log.Printf("Creating secret in namespace %s\n", namespace)
	home := os.Getenv("HOME")
	keyFileName := "ixia-c-automation.json"
	keyFilePath := fmt.Sprintf("%s/.ixops/%s", home, keyFileName)
	data, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		log.Printf("Failed to read key file error - %v\n", err)
		return err
	}

	args := []string{"create", "secret", "-n", namespace, "docker-registry", "kne-pull-secret"}
	args = append(args, "--docker-server=us-central1-docker.pkg.dev")
	args = append(args, "--docker-username=_json_key")
	args = append(args, fmt.Sprintf("--docker-password=\"%s\"", data))
	args = append(args, fmt.Sprintf("--docker-email=%s", setup.GCloudEmail))
	_, err = utils.ExecCmd("kubectl", args...)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Printf("Failed to create kne secret error - %v\n", err)
		return err
	}
	args[5] = "ixia-pull-secret"
	_, err = utils.ExecCmd("kubectl", args...)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Printf("Failed to create ixia secret error - %v\n", err)
		return err
	}

	log.Printf("Secrets created successfully")
	return nil
}

func CreateTopologyWithFile(topologyFile string) error {
	errorString := ""
	userHome, err := os.UserHomeDir()
	if err != nil {
		errorString = fmt.Sprintf("failed to get home: %v\n", err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	topologyFileBasePath, err := utils.FileRelative(topologyFile)
	if err != nil {
		errorString = fmt.Sprintf("failed to get relative path: %v\n", err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	if setup.ClusterTypeGC {
		err = createSecrets(topologyFile)
		if err != nil {
			return err
		}
	}

	kubeCfgLoc := filepath.Join(userHome, ".kube", "config")
	topologyParams := topo1.TopologyParams{
		TopoName:       topologyFile,
		Kubecfg:        kubeCfgLoc,
		TopoNewOptions: []topo1.Option{topo1.WithBasePath(topologyFileBasePath)},
		Timeout:        0,
		DryRun:         false,
	}

	err = topo1.CreateTopology(context.Background(), topologyParams)
	if err != nil {
		errorString = fmt.Sprintf("failed to create topology: %v\n", err)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	namespace, err := getTopologyNamespace(topologyFile)
	if err != nil {
		return err
	}

	log.Printf("Waiting for pods to be ready...")
	kubeClient, err := utils.NewK8sClient()
	if err != nil {
		errorString = fmt.Sprintf("failed to create k8s client: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	err = utils.WaitFor(func() (bool, error) { return kubeClient.AllPodsAreReady(namespace) }, nil)
	if err != nil {
		errorString = fmt.Sprintf("%s pods are not ready: %v", namespace, err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	return nil
}
