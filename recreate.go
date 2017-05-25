package main

import (
  "fmt"
  "os"
  "strings"

  "github.com/fsouza/go-dockerclient"
  //"github.com/tonnerre/golang-pretty"
)

func checkError(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(0)
  }
}

func parseImageName(imageName string) (repository string, tag string) {
  sepIndex := strings.LastIndex(imageName, ":")

  if sepIndex > -1 {
    repository := imageName[:sepIndex]
    tag := imageName[(sepIndex+1):]

    return repository, tag
  } else {
    return imageName, "latest"
  }
}

func main() {
  if len(os.Args) < 2 {
    fmt.Printf("Usage: %s [-p] id\n", os.Args[0])
    os.Exit(0)
  }

  endpoint := "unix:///var/run/docker.sock"
  client, _ := docker.NewClient(endpoint)

  pullImage := os.Args[1] == "-p"
  containerId := os.Args[len(os.Args) - 1]

  oldContainer, err := client.InspectContainer(containerId)
  checkError(err)

  // TODO delete _new if an error occures

  repository, tag := parseImageName(oldContainer.Config.Image)
  fmt.Printf("Image: %s:%s\n", repository, tag)

  if pullImage {
    fmt.Print("Pulling image...\n")

    err = client.PullImage(docker.PullImageOptions{
      Repository: repository,
      Tag: tag }, docker.AuthConfiguration{})

    checkError(err)
  }

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
  checkError(err)

  // rename
  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: oldContainer.ID,
    Name: name + "_old" })
  checkError(err)

  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: newContainer.ID,
    Name: name})
  checkError(err)

  if oldContainer.State.Running {
    fmt.Printf("Stopping old container\n")
    err = client.StopContainer(oldContainer.ID, 10)
    checkError(err)

    fmt.Printf("Starting new container\n")
    err = client.StartContainer(newContainer.ID, newContainer.HostConfig)
    checkError(err)
  }

  // TODO fallback to old container if error occured
  // TODO add option to remove old container on sucsess

  fmt.Printf(
    "Migrated from %s to %s\n",
    oldContainer.ID[:4],
    newContainer.ID[:4])

  fmt.Println("Done")
}
