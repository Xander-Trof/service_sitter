package dockercomands

import (
    "strings"

    "github.com/moby/moby/api/types/container"
)

func FindContainerByName(containerList []container.Summary, targetName string) *container.Summary {
    for _, container := range containerList {
        for _, name := range container.Names {
            // Убираем префикс '/' из имени контейнера
            cleanName := strings.TrimPrefix(name, "/")
            if cleanName == targetName {
                return &container
            }
        }
    }
    return nil
}

func GetContainerNames(containerList []container.Summary) []string {
    names := make([]string, 0, len(containerList))
    for _, container := range containerList {
        for _, name := range container.Names {
            cleanName := strings.TrimPrefix(name, "/")
            names = append(names, cleanName)
        }
    }
    return names
}
