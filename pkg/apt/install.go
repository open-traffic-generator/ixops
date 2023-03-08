package apt

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func (a *aptInstaller) CheckCmd(cmd string) error {
	log.Trace().Str("cmd", cmd).Msg("Checking if command executes successfully")
	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("check for cmd '%s' failed: %v", cmd, a.executor.StderrLines())
	}

	return nil
}

func (a *aptInstaller) InstallPackage(pkgName string, pkgCheckCmd string) error {
	if a.CheckCmd(pkgCheckCmd) == nil {
		log.Trace().Str("pkgName", pkgName).Msg("Skipping installation")
		return nil
	}

	log.Info().Str("pkgName", pkgName).Msg("Installing package")
	if a.installSudo {
		a.installSudo = false
		if err := a.InstallSudo(); err != nil {
			return fmt.Errorf("could not install package %s: %v", pkgName, err)
		}
		a.useSudo = true
	}

	if err := a.Install([]string{pkgName}); err != nil {
		return fmt.Errorf("could not install package %s: %v", pkgName, err)
	}

	if err := a.CheckCmd(pkgCheckCmd); err != nil {
		return fmt.Errorf("post install check for package %s failed: %v", pkgName, err)
	}

	log.Info().Str("pkgName", pkgName).Msg("Successfully installed package")
	return nil
}

func (a *aptInstaller) InstallCurl() error {
	return a.InstallPackage("curl", "curl --version")
}

func (a *aptInstaller) InstallGit() error {
	return a.InstallPackage("git", "git version")
}

func (a *aptInstaller) InstallVim() error {
	return a.InstallPackage("vim", "dpkg -s vim")
}

func (a *aptInstaller) InstallLsbRelease() error {
	return a.InstallPackage("lsb_release", "lsb_release -v")
}

func (a *aptInstaller) InstallGnupg() error {
	return a.InstallPackage("gnupg", "gpg -k")
}

func (a *aptInstaller) InstallCaCertificates() error {
	return a.InstallPackage("ca-certificates", "dpkg -s ca-certificates")
}

func (a *aptInstaller) InstallSudo() error {
	return a.InstallPackage("sudo", "sudo -V")
}

func (a *aptInstaller) AddUserToDockerGroup() error {
	if a.CheckCmd("docker ps -a") == nil {
		log.Trace().Msg("Skipping addition of current user to docker group")
		return nil
	}

	log.Trace().Msg("Adding current user to docker group")
	if err := a.executor.Clear().BashExec("sudo groupadd docker").Err(); err != nil {
		log.Trace().Strs("stderr", a.executor.StderrLines()).Msg("Skipped creation of docker group")
	}
	if err := a.executor.Clear().BashExec("sudo usermod -aG docker $USER").Err(); err != nil {
		return fmt.Errorf("could not add current user to group docker: %v", a.executor.StderrLines())
	}

	log.Info().Msg("Please logout, login again and re-execute previous command")
	os.Exit(0)
	return nil
}

func (a *aptInstaller) InstallDocker() error {
	checkCmd := "docker -v"
	gpgRemote := "https://download.docker.com/linux/ubuntu/gpg"
	gpgLocal := "/usr/share/keyrings/docker-archive-keyring.gpg"
	repoRemote := "https://download.docker.com/linux/ubuntu"
	repoLocal := "/etc/apt/sources.list.d/docker.list"

	if a.CheckCmd(checkCmd) == nil {
		log.Trace().Msg("Skipping docker installation")
		if err := a.AddUserToDockerGroup(); err != nil {
			return err
		}
		return nil
	}

	log.Trace().Msg("Removing existing docker components")
	if err := a.UninstallPackages([]string{"docker", "docker-engine", "docker.io", "containerd", "runc"}); err != nil {
		log.Trace().Msg("Some docker components were already uninstalled")
	}

	log.Trace().Msg("Setting up gpg for docker")
	cmd := fmt.Sprintf("curl -kfsSL %s | sudo gpg --batch --yes --dearmor -o %s", gpgRemote, gpgLocal)
	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("Could not setup gpg for docker: %v", a.executor.StderrLines())
	}

	log.Trace().Msg("Setting up repository for docker")
	cmd = fmt.Sprintf(
		"echo \"deb [arch=amd64 signed-by=%s] %s $(lsb_release -cs) stable\" | sudo tee %s",
		gpgLocal, repoRemote, repoLocal,
	)
	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("Could not setup repo for docker: %v", a.executor.StderrLines())
	}

	log.Trace().Msg("Installing new docker components")
	if err := a.Install([]string{"docker-ce", "docker-ce-cli", "containerd.io"}); err != nil {
		return fmt.Errorf("could not new docker components: %v", err)
	}

	if err := a.CheckCmd(checkCmd); err != nil {
		return fmt.Errorf("post install check for docker failed: %v", err)
	}

	if err := a.AddUserToDockerGroup(); err != nil {
		return err
	}

	return nil
}
