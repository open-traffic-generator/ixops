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
