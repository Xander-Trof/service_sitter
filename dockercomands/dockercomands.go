package dockercomands

import (
    "context"

	"github.com/moby/moby/client"
)

func DockerPS() []string {
    cli, err := client.New(client.FromEnv)
    if err != nil {
        panic(err)
    }
    defer cli.Close()

    
    containersListResult, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    containers := containersListResult.Items
    
    containerNames := make([]string, len(containers))
    for i, container := range containers {
        containerNames[i] = container.Names[0]
    }

    return containerNames
}
