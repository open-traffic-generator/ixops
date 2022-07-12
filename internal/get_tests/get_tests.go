package get_tests

import (
	"github.com/open-traffic-generator/ixops/pkg/configs"
)

func GetTests(c configs.AppConfig) error {
	// errorString := ""

	// home, _ := utils.GetHomeDirectory()
	// ixiaCclientTestHome := utils.ReturnPath([]string{home, c.IxOpsHome,
	// 	c.IxiaC.Home, c.IxiaC.TestClient.Home})

	// ns := "ixia-c-tests"
	// podName := "ixia-c-test-client"

	// testClientYaml, err := mkTcPodYaml(home, c)
	// if err != nil {
	// 	return err
	// }

	// _, err = utils.ExecCmd("kubectl", "create", "namespace", ns)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// _, err = utils.ExecCmd("kubectl", "apply", "-f", testClientYaml)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// getTests(c, ixiaCclientTestHome)

	// log.Printf("Waiting for pods to be ready...")
	// kubeClient, err := utils.NewK8sClient()
	// if err != nil {
	// 	errorString = fmt.Sprintf("failed to create k8s client: %v", err)
	// 	log.Println(err)
	// 	return fmt.Errorf(errorString)
	// }
	// err = utils.WaitFor(func() (bool, error) { return kubeClient.AllPodsAreReady(ns) }, nil)
	// if err != nil {
	// 	errorString = fmt.Sprintf("ixia-c-tests pods are not ready: %v", err)
	// 	log.Println(err)
	// 	return fmt.Errorf(errorString)
	// }

	// _, err = utils.ExecCmd("kubectl", "cp", fmt.Sprintf("%s.tar.gz", ixiaCclientTestHome), fmt.Sprintf("%s/%s:/root/tests.tar.gz", ns, podName))
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// err = setupTc(c, *kubeClient, ns, podName)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	return nil
}

// func getTests(c config.Config, ixiaCclientTestHome string) error {
// 	log.Print("Getting ixia-c test")
// 	home, err := utils.GetHomeDirectory()
// 	if err != nil {
// 		return err
// 	}

// 	ixiaCHome := utils.ReturnPath([]string{home, c.IxOpsHome, c.IxiaC.Home})
// 	utils.ExecCmd("rm", "-rf", ixiaCclientTestHome)
// 	oldPwd, _ := os.Getwd()
// 	os.Chdir(ixiaCHome)
// 	utils.ExecCmd("git", "clone", c.IxiaC.TestClient.Repo.Url)
// 	os.Chdir(ixiaCclientTestHome)
// 	utils.ExecCmd("git", "checkout", c.IxiaC.TestClient.Repo.Commit)
// 	os.Chdir(ixiaCHome)
// 	utils.ExecCmd("rm", "-rf", "tests.tar.gz")
// 	utils.ExecCmd("tar", "czvf", "tests.tar.gz", "tests")
// 	os.Chdir(oldPwd)

// 	return nil
// }

// func mkTcPodYaml(home string, c config.Config) (string, error) {

// 	tcPodConfig := TcPodConfig{
// 		APIVersion: "v1",
// 		Kind:       "Pod",
// 		Metadata: TcPodMetadata{
// 			Name:      "ixia-c-test-client",
// 			Namespace: "ixia-c-tests",
// 		},
// 		Spec: TcPodSpec{
// 			Containers: []Container{},
// 		},
// 	}

// 	tcPodConfig.Spec.Containers = append(tcPodConfig.Spec.Containers, Container{Name: "ubuntu", Image: "ubuntu:22.04", Command: []string{"sleep", "inf"}})

// 	errorString := ""
// 	yamlData, err := yaml.Marshal(&tcPodConfig)
// 	if err != nil {
// 		errorString = fmt.Sprintf("Error while Marshaling. %v", err)
// 		log.Println(err)
// 		return "", fmt.Errorf(errorString)
// 	}
// 	log.Printf("TcPodConfig Config: %s", string(yamlData))

// 	filePath := utils.ReturnPath([]string{home, c.IxOpsHome, c.IxiaC.Home, c.IxiaC.TestClient.Yaml})
// 	fmt.Println(filePath)
// 	err = ioutil.WriteFile(filePath, yamlData, 0666)
// 	if err != nil {
// 		errorString = fmt.Sprintf("Error while wring to %s: %v", filePath, err)
// 		log.Println(err)
// 		return "", fmt.Errorf(errorString)
// 	}

// 	return filePath, nil

// }

// func setupTc(c config.Config, k utils.K8sClient, ns string, podName string) error {
// 	pod, _ := k.GetPod(ns, podName)

// 	command := []string{"/bin/sh", "-c", "apt-get update"}
// 	_, err := k.Exec(pod, command)

// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	command = []string{"/bin/sh", "-c", "apt-get install -y --no-install-recommends sudo build-essential libpcap-dev tar telnet"}
// 	_, err = k.Exec(pod, command)

// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	command = []string{"/bin/sh", "-c", "sudo apt-get install -y --no-install-recommends curl git vim unzip apt-transport-https ca-certificates gnupg lsb-release"}
// 	_, err = k.Exec(pod, command)

// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	getGo := fmt.Sprintf("curl -kL https://dl.google.com/go/go%s.linux-amd64.tar.gz | sudo tar -C /usr/local/ -xzf -", c.GoVersion)

// 	command = []string{"/bin/sh", "-c", getGo}
// 	_, err = k.Exec(pod, command)

// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	s := "'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin'"
// 	command = []string{"/bin/bash", "-c", fmt.Sprintf("echo %s >> ~/.bashrc", s)}
// 	_, err = k.Exec(pod, command)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	command = []string{"/bin/sh", "-c", "cd /root && tar xzvf tests.tar.gz"}
// 	_, err = k.Exec(pod, command)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	command = []string{"/bin/sh", "-c", "rm -rf /root/tests.tar.gz"}
// 	_, err = k.Exec(pod, command)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	command = []string{"/bin/sh", "-c", "cd /root/tests/go/tests && . ~/.bashrc && go test -list . -tags sanity"}
// 	_, err = k.Exec(pod, command)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	return nil
// }
