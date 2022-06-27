package setup

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
)

type Platform struct {
	Os        string
	OsVersion string
}

func checkPlatform() error {
	log.Print("Check Platform")
	gpi, _ := getPlatformInfo() // Need to add Error Handling
	fullVer := strings.Split(gpi.OsVersion, ".")
	mainVer, _ := strconv.ParseInt(fullVer[0], 10, 64)

	comp := strings.Compare(string(gpi.Os), "Ubuntu")
	if mainVer >= 20 && comp == 0 {
		log.Print("Platform Check Passed")
	} else {
		return fmt.Errorf("the tools works only with Ubuntu 20 or above")
	}
	return nil
}

func getPlatformInfo() (Platform, error) {
	gio := Platform{}
	os := ""
	out, err := utils.ExecCmd("grep", "Ubuntu", "/etc/os-release")
	if err != nil {
		return gio, err
	}

	if strings.Contains(out, "Ubuntu") {
		os = "Ubuntu"
	}

	out, err = utils.ExecCmd("grep", "VERSION_ID", "/etc/os-release")
	if err != nil {
		return gio, err
	}
	version := strings.SplitAfter(out, "=")
	osVersion := strings.ReplaceAll(version[1], "\"", "")

	gio = Platform{Os: os, OsVersion: osVersion}
	return gio, err
}

func checkUser(args *[]string) error {
	if len(*args) == 0 {
		return errors.New("setup needs gcloud mail as an argument Ex: ixops cluster setup your.email@example.com")
	} else {
		return nil
	}
}

func setupIxopsHome() error {
	log.Print("Setting up Ixops Home Directory")
	user, err := user.Current()
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	homeDirectory := user.HomeDir
	createDirectory(homeDirectory + "/.ixops")
	return nil
}

func createDirectory(p string) (*os.File, error) {
	if err := os.Mkdir(p, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func createFileAndwrite(file string, bytes []byte) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

	f.Close()
}

func getSysPkgs() error {
	log.Print("Getting Sys Packages to support Ixops")
	_, err := utils.ExecCmd("sudo", "apt-get", "update")
	if err != nil {
		log.Println("sudo apt-get update failed")
		return err
	}
	_, err = utils.ExecCmd("sudo", "apt-get", "install", "-y", "--no-install-recommends", "curl", "git", "vim", "unzip", "apt-transport-https", "ca-certificates", "gnupg", "lsb-release")
	if err != nil {
		log.Println("sudo apt-get update failed")
		return err
	}
	return nil
}

func dockerExists() bool {
	log.Println("Checking for existing Docker....")
	out, _ := utils.ExecCmd("which", "docker")
	if strings.Contains(out, "docker") {
		log.Print("Docker Already Exists")
		return true
	} else {
		log.Print("Docker doesn't Exist, Installing Docker")
		return false
	}
}

func getAndInstallDocker() error {
	log.Print("Installing docker")
	utils.ExecCmd("sudo", "apt-get", "remove", "docker-engine", "docker.io", "containerd", "runc", "2")

	dockergpg()
	dockerList()

	_, err := utils.ExecCmd("sudo", "apt-get", "update")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, err = utils.ExecCmd("sudo", "apt-get", "install", "-y", "docker-ce", "docker-ce-cli", "containerd.io")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	utils.ExecCmd("sudo", "groupadd", "docker")

	currentUser, _ := user.Current()
	_, err = utils.ExecCmd("sudo", "usermod", "-aG", "docker", currentUser.Username)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, err = utils.ExecCmd("sudo", "docker", "version")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func dockerList() {
	lsbRelease, _ := utils.ExecCmd("lsb_release", "-cs")
	lsbRelease = strings.ReplaceAll(lsbRelease, "\n", "")
	cmd1 := exec.Command("echo", "deb", "[arch=amd64", "signed-by=/usr/share/keyrings/docker-archive-keyring.gpg]", "https://download.docker.com/linux/ubuntu", lsbRelease, "stable")
	cmd2 := exec.Command("sudo", "tee", "/etc/apt/sources.list.d/docker.list")

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

func dockergpg() {
	cmd1 := exec.Command("curl", "-kfsSL", "https://download.docker.com/linux/ubuntu/gpg")
	cmd2 := exec.Command("sudo", "gpg", "--batch", "--yes", "--dearmor", "-o", "/usr/share/keyrings/docker-archive-keyring.gpg")

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

func CommonSetup(args *[]string) error {
	err := checkPlatform()
	if err != nil {
		return err
	}

	err = checkUser(args)
	if err != nil {
		return err
	}

	err = setupIxopsHome()
	if err != nil {
		return err
	}

	err = getSysPkgs()
	if err != nil {
		return err
	}

	if !dockerExists() {
		err = getAndInstallDocker()
		if err != nil {
			return err
		}
	}

	_, err = utils.ExecCmd("sudo", "chmod", "666" , "/var/run/docker.sock")
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	
	return nil
}
