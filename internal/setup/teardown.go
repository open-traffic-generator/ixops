package setup

import (
	"fmt"
	// "io/ioutil"
	"log"
	"os"

	// "path/filepath"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
	// "gopkg.in/yaml.v2"
)

func gcDeleteCluster(gcUser string) error {
	gcClusterName := fmt.Sprintf("%s.k8s.local", gcUser)
	gcStoreName := fmt.Sprintf("gs://%s-kops-store/", gcUser)

	//log.Printf("gCloud check if cluster %s (store %s) exists\n", gcClusterName, gcStoreName)
	exists, err := kopsClusterExists(gcClusterName, gcStoreName)
	if err != nil {
		return err
	} else if !exists {
		log.Printf("gCloud cluster %s not found\n", gcClusterName)
		return nil
	}

	log.Printf("gCloud cluster found.\nDeleting cluster %s (store %s)\n", gcClusterName, gcStoreName)
	_, err = utils.ExecCmd("kops", "delete", "cluster", fmt.Sprintf("--name=%s", gcClusterName), fmt.Sprintf("--state=%s", gcStoreName), "--yes")
	if err != nil {
		log.Printf("Failed to delete gCloud cluster error - %v\n", err)
		return err
	}
	return nil
}

func gcDeleteStore(gcUser string) error {
	home := os.Getenv("HOME")
	ixOpsHome := fmt.Sprintf("%s/.ixops", home)
	gUtil := fmt.Sprintf("%s/google-cloud-sdk/bin/gsutil", ixOpsHome)
	gcStoreName := fmt.Sprintf("gs://%s-kops-store/", gcUser)

	log.Printf("gCloud check if store exists\n")
	out, err := utils.ExecCmd(gUtil, "ls")
	if err != nil {
		log.Printf("Failed to check gCloud store error - %v\n", err)
		return err
	} else if !strings.Contains(out, gcStoreName) {
		log.Printf("gCloud store %s not found\n", gcStoreName)
		return nil
	}

	log.Printf("Deleting gCloud store %s\n", gcStoreName)
	_, err = utils.ExecCmd(gUtil, "-m", "rm", "-r", gcStoreName)
	if err != nil {
		log.Printf("failed - %v\n", err)
		return err
	}
	return nil
}

func TeardownCluster(gcCluster bool) error {
	var err error
	if gcCluster {
		err = gcDeleteCluster(GCloudUser)
		if err != nil {
			return err
		}

		err = gcDeleteStore(GCloudUser)
		if err != nil {
			return err
		}

		return nil
	}
	return nil
}
