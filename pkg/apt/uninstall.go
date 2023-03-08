package apt

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (a *aptInstaller) UninstallPackages(pkgNames []string) error {
	log.Info().Strs("pkgNames", pkgNames).Msg("Unstalling packages")
	if a.installSudo {
		a.installSudo = false
		if err := a.InstallSudo(); err != nil {
			return fmt.Errorf("could no uninstall packages %v: %v", pkgNames, err)
		}
		a.useSudo = true
	}

	if err := a.Uninstall(pkgNames); err != nil {
		return fmt.Errorf("could no uninstall packages %v: %v", pkgNames, err)
	}

	log.Info().Strs("pkgNames", pkgNames).Msg("Successfully uninstalled packages")
	return nil
}

func (a *aptInstaller) UninstallGo() error {
	log.Trace().Msg("Uninstalling existing Go")

	cmd := "rm -rf ${HOME}/.local/go ${HOME}/go"
	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("could not uninstall existing Go: %v", a.executor.StderrLines())
	}

	return nil
}
