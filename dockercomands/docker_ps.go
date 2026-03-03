package dockercomands

import (
    "context"

	"github.com/moby/moby/client"
)


func DockerPS() client.ContainerListResult {
    cli, err := client.New(client.FromEnv)
    if err != nil {
        panic(err)
    }
    defer cli.Close()

    containersListResult, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    return containersListResult
}
