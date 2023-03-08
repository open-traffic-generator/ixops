package api

type PackageInstaller interface {
	Update() error
	Install([]string) error
	Uninstall([]string) error

	InstallSudo() error
	InstallCurl() error
	InstallGit() error
	InstallVim() error
	InstallLsbRelease() error
	InstallGnupg() error
	InstallCaCertificates() error
	InstallDocker() error
	InstallGo(string) error

	UninstallGo() error
}
