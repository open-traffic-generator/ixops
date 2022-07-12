package interfaces

import (
	"fmt"

	"github.com/open-traffic-generator/ixops/pkg/exec"
	"github.com/rs/zerolog/log"
)

func CreateVethPair(ifcA string, ifcZ string) error {
	log.Debug().Str("ifcA", ifcA).Str("ifcZ", ifcZ).Msg("Creating veth pair")
	commands := []string{
		fmt.Sprintf("sudo ip link add %s type veth peer name %s", ifcA, ifcZ),
		fmt.Sprintf("sudo ip link set %s up", ifcA),
		fmt.Sprintf("sudo ip link set %s up", ifcZ),
	}
	if err := exec.ExecBashCmd(commands); err != nil {
		return fmt.Errorf("could not create veth pair: %v", err)
	}
	log.Debug().Msg("Created veth pair")
	return nil
}

func DeleteVethPair(ifcA string, ifcZ string) error {
	log.Debug().Str("ifcA", ifcA).Str("ifcZ", ifcZ).Msg("Deleting veth pair")
	commands := []string{
		fmt.Sprintf("sudo ip link delete %s", ifcA),
	}
	if err := exec.ExecBashCmd(commands); err != nil {
		return fmt.Errorf("could not delete veth pair: %v", err)
	}
	log.Debug().Msg("Deleted veth pair")
	return nil
}
