package recreate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func generateContainerNames(container *docker.Container) (
	temporaryName string,
	recentName string,
) {
	now := int(time.Now().Unix())
	then := now - 1

	name := container.Name
	temporaryName = name + "_" + strconv.Itoa(now)
	recentName = name + "_" + strconv.Itoa(then)

	return temporaryName, recentName
}

func cloneContainerLinks(container *docker.Container) (
	clonedLinks []string,
	err error,
) {
	links := container.HostConfig.Links

	for i := range links {
		parts := strings.SplitN(links[i], ":", 2)

		if len(parts) != 2 {
			// TODO make function and add better error return
			return nil, fmt.Errorf("Unable to parse link %s", links[i])
		}

		containerName := strings.TrimPrefix(parts[0], "/")
		aliasParts := strings.Split(parts[1], "/")
		alias := aliasParts[len(aliasParts)-1]

		links[i] = fmt.Sprintf("%s:%s", containerName, alias)
	}

	return links, nil
}

func cloneContainerOptions(
	container *docker.Container,
	imageURL string,
	containerName string,
) (
	options docker.CreateContainerOptions,
	err error,
) {
	options.Name = containerName
	options.Config = container.Config
	options.Config.Image = imageURL
	options.HostConfig = container.HostConfig
	options.HostConfig.VolumesFrom = []string{container.ID}

	links, err := cloneContainerLinks(container)
	options.HostConfig.Links = links

	return options, err
}

func generateEnvMap(envArray []string) map[string]string {
	envMap := make(map[string]string)

	for _, env := range envArray {
		parts := strings.SplitN(env, "=", 2)

		if len(parts) != 2 {
			continue
		}

		envMap[parts[0]] = parts[1]
	}

	return envMap
}

func mergeContainerEnv(options docker.CreateContainerOptions, envMap map[string]string) []string {
	mergedMap := generateEnvMap(options.Config.Env)
	var mergedArray []string

	for k, v := range envMap {
		mergedMap[k] = v
	}

	for k, v := range mergedMap {
		mergedArray = append(mergedArray, k+"="+v)
	}

	return mergedArray
}
