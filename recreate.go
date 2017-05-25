package main

import (
  "fmt"
  "os"
  "strings"

  "github.com/fsouza/go-dockerclient"
  //"github.com/tonnerre/golang-pretty"
)

func check(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(0)
  }
}

func main() {
  if len(os.Args) < 2 {
    fmt.Printf("Usage: %s id [tag]\n", os.Args[0])
    os.Exit(0)
  }

  endpoint := "unix:///var/run/docker.sock"
  client, _ := docker.NewClient(endpoint)

  containerId := os.Args[1]
  // imageTag := os.Args[2]

  oldContainer, err := client.InspectContainer(containerId)
  check(err)

  //TODO delete _new if an error occures

  // TODO check to make sure rebuilt container is using updated image
  // TODO pull image beforehand
  // TODO parse image tag (none: latest?)
  fmt.Printf("Image: %s\n", oldContainer.Config.Image)

  // TODO handle image tags/labels?

  name := oldContainer.Name
  temporaryName := name + "_new"

  // TODO possibility to add/change environment variables
  var options docker.CreateContainerOptions
  options.Name = temporaryName
  options.Config = oldContainer.Config
  options.HostConfig = oldContainer.HostConfig
  options.HostConfig.VolumesFrom = []string{oldContainer.ID}

  links := oldContainer.HostConfig.Links

  for i := range links {
    parts := strings.SplitN(links[i], ":", 2)
    if len(parts) != 2 {
      fmt.Println("Unable to parse link ", links[i])
      // TODO make function and add better error return
      return
    }

    containerName := strings.TrimPrefix(parts[0], "/")
    aliasParts := strings.Split(parts[1], "/")
    alias := aliasParts[len(aliasParts)-1]
    links[i] = fmt.Sprintf("%s:%s", containerName, alias)
  }
  options.HostConfig.Links = links


  fmt.Println("Creating...")
  newContainer, err := client.CreateContainer(options)
  check(err)

  // rename
  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: oldContainer.ID,
    Name: name + "_old" })
  check(err)

  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: newContainer.ID,
    Name: name})
  check(err)

  if oldContainer.State.Running {
    fmt.Printf("Stopping old container\n")
    err = client.StopContainer(oldContainer.ID, 10)
    check(err)

    fmt.Printf("Starting new container\n")
    err = client.StartContainer(newContainer.ID, newContainer.HostConfig)
    check(err)
  }

  // TODO fallback to old container if error occured
  // TODO add option to remove old container on sucsess

  fmt.Printf("Migrated from %s to %s\n", oldContainer.ID, newContainer.ID)

  fmt.Println("Done")
}
