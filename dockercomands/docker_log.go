package dockercomands

import (
    "context"
    "log"

	"github.com/moby/moby/client"
)

func DockerLogs(serviceName string) client.ContainerLogsResult {
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
        return nil
    }

    res, err := cli.ContainerLogs(context.Background(), container.ID, client.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    // Не закрываем здесь — закрытие происходит в getLogs после чтения

    return res
}