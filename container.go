package main

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
  repository string,
  tagName string,
  containerName string,
) (
  options docker.CreateContainerOptions,
  err error,
) {
  options.Name = containerName
  options.Config = container.Config
  options.Config.Image = repository + ":" + tagName
  options.HostConfig = container.HostConfig
  options.HostConfig.VolumesFrom = []string{container.ID}

  links, err := cloneContainerLinks(container)
  options.HostConfig.Links = links

  return options, err
}
