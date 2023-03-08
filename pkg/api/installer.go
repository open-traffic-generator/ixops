package api

type PackageInstaller interface {
	Update() error
	Install([]string) error

	InstallCurl() error
	InstallGit() error
	InstallVim() error
	InstallLsbRelease() error
	InstallGnupg() error
	InstallCaCertificates() error
	InstallSudo() error
}
