package get_images

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
	"gopkg.in/yaml.v2"
)

func DownloadConfigMap(filepath string, url string) error {
	log.Print("Downloading ConfigMap ixia-configmap.yaml")
	error := utils.ReadUrlAndWrite(filepath, url)
	if error != nil {
		log.Print(error)
		return error
	}
	return nil
}

func fromYaml() *ConfigMap {
	configMap := &ConfigMap{}
	data, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Println(err)
	}

	err = yaml.Unmarshal([]byte(data), &configMap)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return configMap
}

func UpdateConfigMapWithGhrc(filepath string) error {
	log.Print("Replacing Image path from Dockerhub to Ghrc")
	configMap := fromYaml()
	configMap.Data.Versions = strings.Replace(configMap.Data.Versions,
		centralDocker, ghrc, -1)

	// Hardcoding for time being
	configMap.Data.Versions = strings.Replace(configMap.Data.Versions,
		"/ixia-c-controller",
		"/licensed/ixia-c-controller", -1)
	configMap.Data.Versions = strings.Replace(configMap.Data.Versions,
		"/ixia-c-protocol-engine",
		"/licensed/ixia-c-protocol-engine", -1)

	configMapVersions = configMap.Data.Versions

	bytes, err := yaml.Marshal(configMap)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	err = ioutil.WriteFile(filepath, bytes, 0644)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	return nil
}

func DownloadAndLoadAllImages() error {
	log.Print("Download All Docker Images and load into Kind Cluster")
	versions := &Version{}
	err := json.Unmarshal([]byte(configMapVersions), &versions)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	dockerLoginGhrc()
	for _, image := range versions.Images {
		utils.ExecCmd("sudo", "docker", "pull", fmt.Sprintf("%s:%s",
			image.Path, image.Tag))
	}

	// Pull DUT image static for now, have to update later
	utils.ExecCmd("sudo", "docker", "pull", "ghcr.io/open-traffic-generator/ceos:4.28.01f")

	for _, image := range versions.Images {
		utils.ExecCmd("kind", "load", "docker-image", fmt.Sprintf("%s:%s",
			image.Path, image.Tag))
	}

	// Pull DUT image static for now, have to update later
	utils.ExecCmd("kind", "load", "docker-image", "ghcr.io/open-traffic-generator/ceos:4.28.01f")

	return nil

}

func dockerLoginGhrc() {
	cmd1 := exec.Command("echo", Pat)
	cmd2 := exec.Command("sudo", "docker", "login", "ghcr.io", "-u", "USERNAME", "--password-stdin")

	// Get the pipe of Stdout from cmd1 and assign it
	// to the Stdin of cmd2.
	pipe, err := cmd1.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd2.Stdin = pipe

	// Start() cmd1, so we don't block on it.
	err = cmd1.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Run Output() on cmd2 to capture the output.
	output, err := cmd2.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(output))
}

func GetImages() error {

	url := fmt.Sprintf("%s/v%s/ixia-configmap.yaml", OtgGit, Ixiacversion)
	err := DownloadConfigMap(FileName, url)
	if err != nil {
		return err
	}
	err = UpdateConfigMapWithGhrc(FileName)
	if err != nil {
		return err
	}
	err = DownloadAndLoadAllImages()
	if err != nil {
		return err
	}
	return nil
}
