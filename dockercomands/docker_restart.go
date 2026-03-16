package dockercomands

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func DockerRestart(serviceName string) (string, error) {
	timeout := 10

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	allContainers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
        All: true,
    })
    if err != nil {
        panic(err)
    }
	restartingContainer := FindContainerByName(allContainers, serviceName)

    if restartingContainer == nil {
        return "", nil
    }

    err = cli.ContainerRestart(context.Background(), restartingContainer.ID, container.StopOptions{
        Timeout: &timeout,
    })
	return restartingContainer.ID, err
}

func DockerRestartAll() []string {
	timeout := 10

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	allContainers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
        All: true,
    })
    if err != nil {
        panic(err)
    }

	allResults := make([]string, 0, len(allContainers))
	for _, RestartingContainer := range allContainers {
		_ = cli.ContainerRestart(context.Background(), RestartingContainer.ID, container.StopOptions{
			Timeout: &timeout,
		})
		allResults = append(allResults, RestartingContainer.ID)
	}
	return allResults
}