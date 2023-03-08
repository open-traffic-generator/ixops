package apt

import "testing"

func TestPackageInstall(t *testing.T) {
	a := NewInstaller()
	if err := a.InstallGit(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
	if err := a.InstallCurl(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
	if err := a.InstallVim(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
	if err := a.InstallLsbRelease(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
	if err := a.InstallGnupg(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
	if err := a.InstallCaCertificates(); err != nil {
		t.Errorf("aptInstaller.InstallGit() error = %v", err)
	}
}
