package recreate

import (
	"github.com/fsouza/go-dockerclient"
)

// Recreation describes a recreation step
type Recreation struct {
	PreviousContainerID string `json:"previousContainerID"`
	NewContainerID      string `json:"newContainerID"`
}

// DockerOptions describe additional options for pulling and creating the container
type DockerOptions struct {
	PullImage       bool
	DeleteContainer bool
	Registries      []RegistryConf
}

// ContainerOptions describe additional options applied to the container
type ContainerOptions struct {
}

// Recreate a container within the default Docker environment
func Recreate(
	containerID string,
	imageTag string,
	containerOptions *ContainerOptions,
	dockerOptions *DockerOptions) (
	recreation *Recreation,
	err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	recreation, err = RecreateWithClient(containerID, imageTag, containerOptions, dockerOptions, client)
	if err != nil {
		return nil, err
	}

	return recreation, nil
}

// RecreateWithEndpoint a container on a specified endpoint
func RecreateWithEndpoint(
	containerID string,
	imageTag string,
	containerOptions *ContainerOptions,
	dockerOptions *DockerOptions,
	endpoint string) (
	recreation *Recreation,
	err error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	recreation, err = RecreateWithClient(containerID, imageTag, containerOptions, dockerOptions, client)
	if err != nil {
		return nil, err
	}

	return recreation, nil
}

// RecreateWithClient recreates a container with a given Docker client
func RecreateWithClient(
	containerID string,
	imageTag string,
	containerOptions *ContainerOptions,
	dockerOptions *DockerOptions,
	client *docker.Client) (recreation *Recreation, err error) {
	previousContainer, err := client.InspectContainer(containerID)
	if err != nil {
		return nil, err
	}

	imageSpec := parseImageName(previousContainer.Config.Image)

	if imageTag != "" {
		imageSpec.tag = imageTag
	}

	if dockerOptions.PullImage {
		auth := findRegistry(dockerOptions.Registries, imageSpec.registry)
		pullOpts := docker.PullImageOptions{
			Repository: imageSpec.repository,
			Tag:        imageSpec.tag}

		err = client.PullImage(
			pullOpts,
			auth)

		if err != nil {
			return nil, err
		}
	}

	temporaryName, recentName := generateContainerNames(previousContainer)

	cloneOptions, err := cloneContainerOptions(
		previousContainer,
		imageSpec.repository,
		imageSpec.tag,
		temporaryName)

	if err != nil {
		return nil, err
	}

	newContainer, err := client.CreateContainer(cloneOptions)

	if err != nil {
		return nil, err
	}

	err = client.RenameContainer(docker.RenameContainerOptions{
		ID:   previousContainer.ID,
		Name: recentName})

	if err != nil {
		return nil, err
	}

	err = client.RenameContainer(docker.RenameContainerOptions{
		ID:   newContainer.ID,
		Name: previousContainer.Name})

	if err != nil {
		return nil, err
	}

	if previousContainer.State.Running {
		err = client.StopContainer(previousContainer.ID, 10)

		if err != nil {
			return nil, err
		}

		err = client.StartContainer(newContainer.ID, newContainer.HostConfig)

		if err != nil {
			return nil, err
		}
	}

	if dockerOptions.DeleteContainer {
		err = client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            previousContainer.ID,
			RemoveVolumes: false})

		if err != nil {
			return nil, err
		}
	}

	return &Recreation{
			PreviousContainerID: previousContainer.ID,
			NewContainerID:      newContainer.ID},
		nil
}
