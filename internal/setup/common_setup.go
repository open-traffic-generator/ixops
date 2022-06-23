package setup

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/open-traffic-generator/ixops/internal/utils"
)

type Platform struct {
	Os        string
	OsVersion string
}

func checkPlatform() error {
	gpi, _ := getPlatformInfo() // Need to add Error Handling
	res := strings.Split(gpi.OsVersion, " ")

	fullVer := strings.Split(res[1], ".")
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
	out, err := utils.ExecCmd("grep", "Ubuntu", "/etc/os-release")

	out = strings.TrimSpace(out)
	var lines []string = regexp.MustCompile("\r?\n").Split(out, -1)
	res1 := strings.SplitAfter(lines[0], "=")
	res2 := strings.SplitAfter(lines[1], "=")

	os := strings.TrimSpace(res1[1])
	osversion := strings.TrimSpace(res2[1])

	os = strings.ReplaceAll(os, "\"", "")
	gio := Platform{Os: os, OsVersion: osversion}
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

func getAndInstallDocker() error {
	log.Print("Installing docker")
	_, err := utils.ExecCmd("sudo", "apt-get", "remove", "docker-engine", "docker.io")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// Need to check this
	// gpgOut, _ := utils.ExecCmd("curl", "-kfsSL", "https://download.docker.com/linux/ubuntu/gpg")
	// gpgOut = strings.TrimRight(gpgOut, "\n")
	// utils.ExecCmd("sudo", "gpg", "--batch", "--yes", "--dearmor", "-o", "/usr/share/keyrings/docker-archive-keyring.gpg")
	// utils.ExecCmd("sudo", "chmod", "777", "/usr/share/keyrings/docker-archive-keyring.gpg")
	// createFileAndwrite("/usr/share/keyrings/docker-archive-keyring.gpg", []byte(gpgOut))

	lsbRelease, _ := utils.ExecCmd("lsb_release", "-cs")
	lsbRelease = strings.ReplaceAll(lsbRelease, "\n", "")
	out, err := utils.ExecCmd("echo", "deb", "[arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg]", "https://download.docker.com/linux/ubuntu", lsbRelease, "stable")
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	_, err = utils.ExecCmd("sudo", "chmod", "777", "/etc/apt/sources.list.d/docker.list")
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	createFileAndwrite("/etc/apt/sources.list.d/docker.list", []byte(out))

	_, err = utils.ExecCmd("sudo", "apt-get", "update")
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

	err = getAndInstallDocker()
	if err != nil {
		return err
	}
	return nil
}
