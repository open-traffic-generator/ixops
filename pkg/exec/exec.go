package exec

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

func ExecBashCmd(commands []string) error {
	log.Trace().Strs("commands", commands).Msg("Executing bash commands")
	cmd := exec.Command("bash", "-c", strings.Join(commands, "&&"))

	o, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("could not start execution: %v", err)
	}

	log.Trace().Str("output", string(o)).Msg("Executed bash command")
	return nil
}
