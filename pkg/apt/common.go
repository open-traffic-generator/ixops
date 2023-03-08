package apt

import (
	"fmt"
	"strings"

	"github.com/open-traffic-generator/ixops/pkg/api"
	"github.com/open-traffic-generator/ixops/pkg/exec"
	"github.com/rs/zerolog/log"
)

type aptInstaller struct {
	update      bool
	useSudo     bool
	installSudo bool
	executor    exec.Executor
}

func NewInstaller() api.PackageInstaller {
	return &aptInstaller{
		update:      true,
		useSudo:     false,
		installSudo: true,
		executor:    exec.NewExecutor(),
	}
}

func (a *aptInstaller) SetUpdate(update bool) api.PackageInstaller {
	a.update = update
	return a
}

func (a *aptInstaller) Install(packages []string) error {
	log.Info().Strs("packages", packages).Msg("Performing apt-get install")

	if err := a.Update(); err != nil {
		return fmt.Errorf("could not install packages %v: %v", packages, err)
	}

	cmd := fmt.Sprintf(
		"DEBIAN_FRONTEND=noninteractive apt-get install -yq --no-install-recommends %s",
		strings.Join(packages, " "),
	)
	if a.useSudo {
		cmd = "sudo " + cmd
	}

	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("Could not perform apt-get install: %v", a.executor.StderrLines())
	}

	return nil
}

func (a *aptInstaller) Update() error {
	if !a.update {
		log.Trace().Msg("Skipping apt-get update")
		return nil
	}
	log.Info().Msg("Performing apt-get update")

	cmd := "DEBIAN_FRONTEND=noninteractive apt-get update -yq --no-install-recommends"
	if a.useSudo {
		cmd = "sudo " + cmd
	}

	if err := a.executor.Clear().BashExec(cmd).Err(); err != nil {
		return fmt.Errorf("Could not perform apt-get update: %v", a.executor.StderrLines())
	}

	log.Info().Msg("Successfully performed apt-get update")
	a.update = false
	return nil
}
