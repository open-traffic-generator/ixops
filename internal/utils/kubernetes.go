package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	core_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sClient struct {
	config    *rest.Config
	clientSet *kubernetes.Clientset
}

func NewK8sClient() (*K8sClient, error) {
	log.Printf("Creating k8s client...")
	var config *rest.Config
	var err error

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get home: %v\n", err)
		return nil, fmt.Errorf(fmt.Sprintf("failed to get home: %v\n", err))
	}

	kubeCfgLoc := filepath.Join(userHome, ".kube", "config")
	log.Printf("Using .kube/config...")
	config, err = clientcmd.BuildConfigFromFlags("", kubeCfgLoc)
	if err != nil {
		return nil, fmt.Errorf("error in building config from %s: %s", kubeCfgLoc, err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error in creating clientset: %s", err.Error())
	}
	log.Println("Successfully created k8s client !")

	return &K8sClient{
		config:    config,
		clientSet: clientset,
	}, nil
}

func (c *K8sClient) AllPodsAreReady(namespace string) (bool, error) {
	pods, err := c.clientSet.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return false, err
	}
	if len(pods.Items) == 0 {
		return false, nil
	}
	for _, pod := range pods.Items {
		for _, condition := range pod.Status.Conditions {
			log.Printf("waiting for pod %s in namespace %s to be ready....", pod.Name, namespace)
			if condition.Status != "True" {
				return false, nil
			}
			log.Printf("pod %s in namespace %s: condition: %s, is ready", pod.Name, namespace, condition.Type)
		}
	}
	return true, nil
}

func (c *K8sClient) GetPod(namespace string, podname string) (p *core_v1.Pod, err error) {

	pod, err := c.clientSet.CoreV1().Pods(namespace).Get(context.TODO(), podname, v1.GetOptions{})
	if err != nil {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", podname, namespace)
	}
	return pod, nil

}

func (c *K8sClient) Exec(p *core_v1.Pod, command []string) (string, error) {

	attachOptions := &core_v1.PodExecOptions{
		Container: p.Spec.Containers[0].Name,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Command:   command,
	}

	request := c.clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(p.Name).
		Namespace(p.Namespace).
		SubResource("exec").
		VersionedParams(attachOptions, scheme.ParameterCodec)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	streamOptions := remotecommand.StreamOptions{
		Stdout: stdout,
		Stderr: stderr,
	}

	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", request.URL())
	if err != nil {
		fmt.Println(exec)
		return "", err
	}

	err = exec.Stream(streamOptions)
	if err != nil {
		result := strings.TrimSpace(stdout.String()) + "\n" + strings.TrimSpace(stderr.String())
		result = strings.TrimSpace(result)
		log.Print(result)
		return "", err
	}

	result := strings.TrimSpace(stdout.String()) + "\n" + strings.TrimSpace(stderr.String())
	result = strings.TrimSpace(result)
	log.Print(result)
	return result, nil
}
