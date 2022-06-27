package run_tests

import (
	"fmt"
	"log"

	"github.com/open-traffic-generator/ixops/internal/utils"
)

func RunTests(args []string) error {

	tcArg := ""
	if len(args) == 0 {
		tcArg = "."
	} else {
		tcArg = "-run " + args[0]
	}

	errorString := ""
	kubeClient, err := utils.NewK8sClient()
	if err != nil {
		errorString = fmt.Sprintf("failed to create k8s client: %v", err)
		log.Println(err)
		return fmt.Errorf(errorString)
	}
	ns := "ixia-c-tests"
	podName := "ixia-c-test-client"

	pod, _ := kubeClient.GetPod(ns, podName)

	command := []string{"/bin/sh", "-c", fmt.Sprintf("cd /root/tests/go/tests && . ~/.bashrc && go test -v %s -tags sanity -timeout 10m", tcArg)}
	_, err = kubeClient.Exec(pod, command)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
