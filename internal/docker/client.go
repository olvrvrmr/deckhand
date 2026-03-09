package docker

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockerclient "github.com/docker/docker/client"
)

type Client struct {
	cli *dockerclient.Client
}

func New() (*Client, error) {
	cli, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client: %w", err)
	}
	return &Client{cli: cli}, nil
}

func (c *Client) GetBackupContainers(ctx context.Context) ([]*ContainerMeta, error) {
	f := filters.NewArgs(filters.Arg("label", LabelEnable+"=true"))
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{Filters: f})
	if err != nil {
		return nil, fmt.Errorf("list containers: %w", err)
	}

	var result []*ContainerMeta
	for _, ct := range containers {
		name := ct.ID[:12]
		if len(ct.Names) > 0 {
			name = strings.TrimPrefix(ct.Names[0], "/")
		}
		if meta := parseLabels(ct.ID, name, ct.Labels); meta != nil {
			result = append(result, meta)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Priority < result[j].Priority
	})

	return result, nil
}

func (c *Client) Stop(ctx context.Context, id string) error {
	timeout := 30
	return c.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

func (c *Client) Start(ctx context.Context, id string) error {
	return c.cli.ContainerStart(ctx, id, container.StartOptions{})
}

func (c *Client) Exec(ctx context.Context, id, command string) error {
	exec, err := c.cli.ContainerExecCreate(ctx, id, container.ExecOptions{
		Cmd:          []string{"sh", "-c", command},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Errorf("exec create: %w", err)
	}
	return c.cli.ContainerExecStart(ctx, exec.ID, container.ExecStartOptions{})
}
