package apt

import (
	"fmt"

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
			return fmt.Errorf("could no install package %s: %v", pkgName, err)
		}
		a.useSudo = true
	}

	if err := a.Install([]string{pkgName}); err != nil {
		return fmt.Errorf("could no install package %s: %v", pkgName, err)
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
