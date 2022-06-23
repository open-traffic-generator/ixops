package utils

import (
	"fmt"
	"log"
	"strings"
)

func GetAllNamespaces() ([]string, error) {
	out, err := ExecCmd("kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		log.Printf("Failed to get all namespaces: %v\n", err)
		return nil, fmt.Errorf(fmt.Sprintf("Failed to get all namespaces: %v", err))
	}
	log.Printf("currently deployed namespaces are: %v\n", strings.Split(out, " "))
	return strings.Split(out, " "), nil
}

func GetPods(namespace string) ([]string, error) {
	out, err := ExecCmd("kubectl", "get", "pods", "-n", namespace, "-o", "jsonpath={.items[*].metadata.name}")
	if err != nil {
		log.Printf("Failed to get all namespaces: %v", err)
		return nil, fmt.Errorf(fmt.Sprintf("Failed to get pods for namespace %s: %v", namespace, err))
	}
	log.Printf("currently deployed pods in namespace %s: %v\n", namespace, strings.Split(out, " "))
	return strings.Split(out, " "), nil
}

func PodReady(pod string, namespace string, waitTime int64) error {
	log.Printf("waiting for pod %s in namespace %s to be ready with in %d seconds\n", pod, namespace, waitTime)
	_, err := ExecCmd("kubectl", "wait", "-n", namespace, fmt.Sprintf("pod/%s", pod), "--for", "condition=ready", "--timeout", fmt.Sprintf("%ds", waitTime))
	if err != nil {
		log.Printf("pod %s in namespace %s is not ready with in %d seconds\n", pod, namespace, waitTime)
		return fmt.Errorf(fmt.Sprintf("pod %s in namespace %s is not ready with in %d seconds\n", pod, namespace, waitTime))
	}
	return nil
}

func WaitForAllPodsToBeReady(namespace string, waitTime int64) error {
	errorString := ""

	actualNamespaces, err := GetAllNamespaces()
	if err != nil {
		return err
	}

	if !Contains(actualNamespaces, namespace) {
		errorString = fmt.Sprintf("%s not found", namespace)
		log.Println(err)
		return fmt.Errorf(errorString)
	}

	actualPods, err := GetPods(namespace)
	if err != nil {
		return err
	}

	for _, pod := range actualPods {
		err := PodReady(pod, namespace, waitTime)
		if err != nil {
			return err
		}
	}
	return nil

}
