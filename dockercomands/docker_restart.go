package dockercomands

import (
	"context"

	"github.com/moby/moby/client"
)

func DockerRestart(serviceName string) (result client.ContainerRestartResult, err error) {
	timeout := 10

	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	allContainers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{
        All: true,
    })
    if err != nil {
        panic(err)
    }
	container := FindContainerByName(allContainers.Items, serviceName)

    if container == nil {
        return client.ContainerRestartResult{}, nil
    }

    result, err = cli.ContainerRestart(context.Background(), container.ID, client.ContainerRestartOptions{
		Timeout: &timeout,
	})
	return result, err
}

func DockerRestartAll() []client.ContainerRestartResult {
	timeout := 10

	cli, err := client.New(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	allContainers, err := cli.ContainerList(context.Background(), client.ContainerListOptions{
        All: true,
    })
    if err != nil {
        panic(err)
    }

	allResults := make([]client.ContainerRestartResult, 0, len(allContainers.Items))
	for _, container := range allContainers.Items {
		result, _ := cli.ContainerRestart(context.Background(), container.ID, client.ContainerRestartOptions{
			Timeout: &timeout,
		})
		allResults = append(allResults, result)
	}
	return allResults
}