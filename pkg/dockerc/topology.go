package dockerc

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/open-traffic-generator/ixops/pkg/interfaces"
	"github.com/rs/zerolog/log"
)

func getIfcPairAndEp(t *configs.Topology, ep *configs.Endpoint, ifcA *string, ifcZ *string) error {
	c := configs.GetAppConfig()
	if len(*c.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint needs to be provided")
	}

	for _, e := range *c.Endpoints {
		if *e.Kind == configs.EndpointKindIxiaC {
			*ep = *e
			break
		}
	}

	if ep == nil {
		return fmt.Errorf("could not find ixia-c endpoint")
	}

	if len(*t.Links) != 1 {
		return fmt.Errorf("exactly one link supported")
	}
	for _, l := range *t.Links {
		pairs := strings.Fields(l)
		if len(pairs) != 2 {
			return fmt.Errorf("exactly two interfaces in link supported")
		}

		splits := strings.Split(pairs[0], ":")
		if len(splits) != 2 {
			return fmt.Errorf("exactly one endpoint name and one interface name supported in a link endpoint")
		}
		epA := splits[0]
		*ifcA = splits[1]

		splits = strings.Split(pairs[1], ":")
		if len(splits) != 2 {
			return fmt.Errorf("exactly one endpoint name and one interface name supported in a link endpoint")
		}
		epZ := splits[0]
		*ifcZ = splits[1]

		if *ep.Name != epA {
			return fmt.Errorf("endpoint %s does not exist", epA)
		}
		if *ep.Name != epZ {
			return fmt.Errorf("endpoint %s does not exist", epZ)
		}
	}

	return nil
}

func CreateTopology(t *configs.Topology) error {
	log.Info().Str("topology", *t.Name).Msg("Creating")
	log.Trace().Interface("topoloy", t).Msg("Creating topology")
	c := configs.GetAppConfig()

	var ifcA string
	var ifcZ string
	var ep configs.Endpoint

	if err := getIfcPairAndEp(t, &ep, &ifcA, &ifcZ); err != nil {
		return err
	}

	if *t.CreateLinks {
		if err := interfaces.CreateVethPair(ifcA, ifcZ); err != nil {
			return err
		}
	}

	dc, err := NewClient((*c.Nodes)[0])
	if err != nil {
		return fmt.Errorf("could not create docker client: %v", err)
	}

	controllerOpts := ControllerOpts{
		Name:        "ixia-c-controller",
		Image:       *ep.IxiaC.Controller,
		AcceptEula:  true,
		Debug:       true,
		HostNetwork: true,
		Port:        443,
	}

	if err := dc.CreateController(&controllerOpts); err != nil {
		return fmt.Errorf("could not create controller: %v", err)
	}

	for i, ifc := range []string{ifcA, ifcZ} {
		trafficEngineOpts := TrafficEngineOpts{
			Name:           fmt.Sprintf("ixia-c-traffic-engine-%s", ifc),
			Image:          *ep.IxiaC.TrafficEngine,
			HugePages:      false,
			CpuPinning:     false,
			WaitForIfc:     true,
			InterfaceNames: []string{ifc},
			HostNetwork:    true,
			Port:           5555 + i,
		}

		if err := dc.CreateTrafficEngine(&trafficEngineOpts); err != nil {
			return fmt.Errorf("could not create traffic-engine: %v", err)
		}
	}

	log.Info().Str("topology", *t.Name).Msg("Successfully created")
	return nil
}

func DeleteTopology(t *configs.Topology) error {
	log.Info().Str("topology", *t.Name).Msg("Deleting")
	log.Trace().Interface("topoloy", t).Msg("Deleting topology")
	c := configs.GetAppConfig()
	var ifcA string
	var ifcZ string
	var ep configs.Endpoint

	if err := getIfcPairAndEp(t, &ep, &ifcA, &ifcZ); err != nil {
		return err
	}

	dc, err := NewClient((*c.Nodes)[0])
	if err != nil {
		return fmt.Errorf("could not create docker client: %v", err)
	}

	ctx := context.Background()
	containers, err := dc.client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return fmt.Errorf("could not list containers: %v", err)
	}

	for _, c := range containers {
		if c.Names[0] == "/ixia-c-controller" {
			if err := dc.DeleteContainer("ixia-c-controller"); err != nil {
				return err
			}
		}
		for _, ifc := range []string{ifcA, ifcZ} {
			cName := fmt.Sprintf("ixia-c-traffic-engine-%s", ifc)
			if c.Names[0] == "/"+cName {
				if err := dc.DeleteContainer(cName); err != nil {
					return err
				}
			}

		}
	}

	if *t.CreateLinks {
		if err := interfaces.DeleteVethPair(ifcA, ifcZ); err != nil {
			return err
		}
	}

	log.Info().Str("topology", *t.Name).Msg("Successfully deleted")
	return nil
}
