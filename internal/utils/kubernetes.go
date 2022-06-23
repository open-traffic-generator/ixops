package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
