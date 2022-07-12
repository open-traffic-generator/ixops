package dockerc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/rs/zerolog/log"
)

type TrafficEngineOpts struct {
	Name           string
	HugePages      bool
	WaitForIfc     bool
	CpuPinning     bool
	Port           int
	InterfaceNames []string
	HostNetwork    bool
	Image          string
}

type ControllerOpts struct {
	AcceptEula  bool
	Name        string
	Image       string
	Debug       bool
	HostNetwork bool
	Port        int
}

func (d *DockerClient) CreateTrafficEngine(opts *TrafficEngineOpts) error {
	if opts == nil {
		return fmt.Errorf("traffic-engine opts not provided")
	}
	if len(opts.InterfaceNames) != 1 {
		return fmt.Errorf("exactly one interface in traffic-engine supported")
	}
	log.Trace().Interface("opts", opts).Msg("Creating traffic-engine container")

	if err := d.PullImage(opts.Image); err != nil {
		return err
	}

	ctx := context.Background()
	envs := []string{}
	if !opts.HugePages {
		envs = append(envs, "OPT_NO_HUGEPAGES=Yes")
	}
	if opts.WaitForIfc {
		envs = append(envs, "WAIT_FOR_IFACE=Yes")
	}
	if !opts.CpuPinning {
		envs = append(envs, "OPT_NO_PINNING=Yes")
	}
	envs = append(envs, fmt.Sprintf("OPT_LISTEN_PORT=%d", opts.Port))
	envs = append(envs, fmt.Sprintf("ARG_IFACE_LIST=virtual@af_packet,%s", opts.InterfaceNames[0]))

	hc := container.HostConfig{
		Privileged: true,
	}
	if opts.HostNetwork {
		hc.NetworkMode = "host"
	}

	r, err := d.client.ContainerCreate(
		ctx, &container.Config{
			Image: opts.Image,
			Env:   envs,
		},
		&hc, nil, nil, opts.Name,
	)
	if err != nil {
		return fmt.Errorf("could not create traffic-engine: %v", err)
	}

	log.Debug().Msg("Created container, now starting it")
	if err := d.client.ContainerStart(ctx, r.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("could not start container: %v", err)
	}

	log.Debug().Msg("Successfully created and started traffic-engine container")
	return nil
}

func (d *DockerClient) CreateController(opts *ControllerOpts) error {
	if opts == nil {
		return fmt.Errorf("controller opts not provided")
	}
	log.Trace().Interface("opts", opts).Msg("Creating controller container")

	if err := d.PullImage(opts.Image); err != nil {
		return err
	}

	ctx := context.Background()
	commands := []string{}
	if opts.AcceptEula {
		commands = append(commands, "--accept-eula")
	}
	if opts.Debug {
		commands = append(commands, "--debug")
	}

	hc := container.HostConfig{}
	if opts.HostNetwork {
		hc.NetworkMode = "host"
	}

	r, err := d.client.ContainerCreate(
		ctx, &container.Config{
			Image: opts.Image,
			Cmd:   commands,
		},
		&hc, nil, nil, opts.Name,
	)
	if err != nil {
		return fmt.Errorf("could not create controller: %v", err)
	}

	log.Debug().Msg("Created container, now starting it")
	if err := d.client.ContainerStart(ctx, r.ID, types.ContainerStartOptions{}); err != nil {
		return fmt.Errorf("could not start container: %v", err)
	}

	log.Debug().Msg("Successfully created and started controller container")
	return nil
}
