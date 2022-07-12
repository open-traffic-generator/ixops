package dockerc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/open-traffic-generator/ixops/pkg/configs"
	"github.com/rs/zerolog/log"
)

type DockerClient struct {
	Host   string
	client *client.Client
}

func NewClient(node *configs.Node) (*DockerClient, error) {
	log.Trace().Interface("node", node).Msg("Creating new docker client")
	dc := &DockerClient{}

	opts := []client.Opt{client.FromEnv, client.WithAPIVersionNegotiation()}
	if node != nil {
		// TODO: allow creating client for remote docker host
		if *node.Host != "localhost" {
			return nil, fmt.Errorf("creating docker client for remote host is currently not supported")
		}
		dc.Host = node.DockerHost()
	}
	c, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not create docker client: %v", err)
	}

	dc.client = c
	log.Trace().Interface("client", dc).Msg("Successfully created docker client")
	return dc, nil
}

func (d *DockerClient) ListImages() error {
	ctx := context.Background()
	summaries, err := d.client.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return fmt.Errorf("could not list images: %v", err)
	}
	log.Trace().Interface("summaries", summaries).Msg("Images")

	return nil
}

func (d *DockerClient) StopContainer(name string) error {
	log.Debug().Str("name", name).Msg("Stopping container")
	ctx := context.Background()

	j, err := d.client.ContainerInspect(ctx, name)
	if err != nil {
		return fmt.Errorf("could not inspect container %s: %v", name, err)
	}

	if j.State.Running {
		log.Debug().Msg("Container running, stopping it")
		if err := d.client.ContainerStop(ctx, j.ID, nil); err != nil {
			return fmt.Errorf("could not stop container %s: %v", j.ID, err)
		}
	}
	log.Debug().Msg("Successfully stopped container")
	return nil
}

func (d *DockerClient) DeleteContainer(name string) error {
	log.Debug().Str("name", name).Msg("Deleting container")
	ctx := context.Background()

	j, err := d.client.ContainerInspect(ctx, name)
	if err != nil {
		return fmt.Errorf("could not inspect container %s: %v", name, err)
	}

	if j.State.Running {
		log.Debug().Msg("Container running, stopping it")
		if err := d.client.ContainerStop(ctx, j.ID, nil); err != nil {
			return fmt.Errorf("could not stop container %s: %v", j.ID, err)
		}
	}

	if err := d.client.ContainerRemove(ctx, j.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("could not delete container %s: %v", j.ID, err)
	}

	return nil
}
