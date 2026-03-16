package dockercomands

import (
    "strings"

    "github.com/docker/docker/api/types"
)

func FindContainerByName(containerList []types.Container, targetName string) *types.Container {
    for i, c := range containerList {
        for _, name := range c.Names {
            // Убираем префикс '/' из имени контейнера
            cleanName := strings.TrimPrefix(name, "/")
            if cleanName == targetName {
                return &containerList[i]
            }
        }
    }
    return nil
}

func GetContainerNames(containerList []types.Container) []string {
    names := make([]string, 0, len(containerList))
    for _, container := range containerList {
        for _, name := range container.Names {
            cleanName := strings.TrimPrefix(name, "/")
            names = append(names, cleanName)
        }
    }
    return names
}
