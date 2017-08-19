package main

import (
  "fmt"
  "os"

  "github.com/fsouza/go-dockerclient"
)

func checkError(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(0)
  }
}

func main() {
  if len(os.Args) < 2 {
    fmt.Printf(`Usage: %s [-p] [-d] id [tag]
  -p Pull image from registry
  -d Delete old container
`, os.Args[0])
    os.Exit(0)
  }

  client, err := docker.NewClientFromEnv()
  checkError(err)

  args, err := parseArgs(os.Args)

  recentContainer, err := client.InspectContainer(args.containerID)
  checkError(err)

  repository, currentTag := parseImageName(recentContainer.Config.Image)

  if args.tagName == "" {
    args.tagName = currentTag
  }

  fmt.Printf("Image: %s:%s\n", repository, args.tagName)

  if args.pullImage {
    fmt.Print("Pulling image...\n")

    err = client.PullImage(docker.PullImageOptions{
      Repository: repository,
      Tag: args.tagName }, docker.AuthConfiguration{})

    checkError(err)
  }

  temporaryName, recentName := generateContainerNames(recentContainer)

  options, err := cloneContainerOptions(
    recentContainer,
    repository,
    args.tagName,
    temporaryName)
  checkError(err)

  fmt.Println("Creating...")
  newContainer, err := client.CreateContainer(options)
  checkError(err)

  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: recentContainer.ID,
    Name: recentName })
  checkError(err)

  err = client.RenameContainer(docker.RenameContainerOptions{
    ID: newContainer.ID,
    Name: recentContainer.Name})
  checkError(err)

  if recentContainer.State.Running {
    fmt.Printf("Stopping old container\n")
    err = client.StopContainer(recentContainer.ID, 10)
    checkError(err)

    fmt.Printf("Starting new container\n")
    err = client.StartContainer(newContainer.ID, newContainer.HostConfig)
    checkError(err)
  }

  if args.deleteContainer {
    fmt.Printf("Deleting old container...\n")

    err = client.RemoveContainer(docker.RemoveContainerOptions{
      ID: recentContainer.ID,
      RemoveVolumes: false })
    checkError(err)
  }

  fmt.Printf(
    "Migrated from %s to %s\n",
    recentContainer.ID[:4],
    newContainer.ID[:4])

  fmt.Println("Done")
}
