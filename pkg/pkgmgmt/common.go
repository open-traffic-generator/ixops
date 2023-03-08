package pkgmgmt

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/pkg/apt"
	"github.com/rs/zerolog/log"
)

type PackageName string
type PackageAction int

const (
	PackageNameAll            PackageName = "all"
	PackageNameSudo           PackageName = "sudo"
	PackageNameVim            PackageName = "vim"
	PackageNameCurl           PackageName = "curl"
	PackageNameGit            PackageName = "git"
	PackageNameLsbRelease     PackageName = "lsb_release"
	PackageNameGnupg          PackageName = "gnupg"
	PackageNameCaCertificates PackageName = "ca-certificates"
	PackageNameDocker         PackageName = "docker"
)

func installAllPackages() error {
	a := apt.NewInstaller()
	if err := a.InstallGit(); err != nil {
		return err
	}
	if err := a.InstallCurl(); err != nil {
		return err
	}
	if err := a.InstallVim(); err != nil {
		return err
	}
	if err := a.InstallLsbRelease(); err != nil {
		return err
	}
	if err := a.InstallGnupg(); err != nil {
		return err
	}
	if err := a.InstallCaCertificates(); err != nil {
		return err
	}
	if err := a.InstallDocker(); err != nil {
		return err
	}
	return nil
}

func InstallPackage(name string, version string) error {
	log.Info().Str("name", name).Str("version", version).Msg("Installing package")

	if PackageName(name) == PackageNameAll {
		return installAllPackages()
	}

	installer := apt.NewInstaller()
	switch PackageName(name) {
	case PackageNameSudo:
		return installer.InstallSudo()
	case PackageNameGit:
		return installer.InstallGit()
	case PackageNameCurl:
		return installer.InstallCurl()
	case PackageNameVim:
		return installer.InstallVim()
	case PackageNameLsbRelease:
		return installer.InstallLsbRelease()
	case PackageNameGnupg:
		return installer.InstallGnupg()
	case PackageNameCaCertificates:
		return installer.InstallCaCertificates()
	case PackageNameDocker:
		return installer.InstallDocker()
	default:
		return fmt.Errorf("unsupported package %s", name)
	}
}

func UninstallPackage(name string) error {
	log.Info().Str("name", name).Msg("Uninstalling package")
	return fmt.Errorf("unimplemented operation")
}
