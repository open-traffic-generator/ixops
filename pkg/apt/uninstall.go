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
