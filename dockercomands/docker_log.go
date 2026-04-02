package dockercomands

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func DockerLogs(serviceName string) io.ReadCloser {
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
	container := FindContainerByName(allContainers, serviceName)

    if container == nil {
        return nil
    }

    res, err := cli.ContainerLogs(context.Background(), container.ID, types.ContainerLogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Since: time.Now().AddDate(0, -1, 0).Format(time.RFC3339),
    })
    if err != nil {
        log.Fatal(err)
    }
    // Не закрываем здесь — закрытие происходит в getLogs после чтения

    return res
}