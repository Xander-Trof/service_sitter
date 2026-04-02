package dockercomands

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)


func DockerPS() []types.Container {
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        panic(err)
    }
    defer cli.Close()

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
    if err != nil {
        panic(err)
    }

    return containers
}
