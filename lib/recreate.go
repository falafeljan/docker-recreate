package recreate

import (
	"github.com/fsouza/go-dockerclient"
)

// Recreation describes a recreation step
type Recreation struct {
	PreviousContainerID string
	NewContainerID      string
}

// Options describe additional options
type Options struct {
	PullImage       bool
	DeleteContainer bool
	Registries      []RegistryConf
}

// Recreate a container with a given Docker client
func Recreate(
	endpoint string,
	containerID string,
	tagName string,
	options *Options) (
	recreation *Recreation,
	err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	recreation, err = RecreateWithClient(client, containerID, tagName, options)
	if err != nil {
		return nil, err
	}

	return recreation, nil
}

// RecreateWithClient recreates a container with a given Docker client
func RecreateWithClient(
	client *docker.Client,
	containerID string,
	tagName string,
	options *Options) (recreation *Recreation, err error) {
	previousContainer, err := client.InspectContainer(containerID)
	if err != nil {
		return nil, err
	}

	imageSpec := parseImageName(previousContainer.Config.Image)

	if tagName != "" {
		imageSpec.tag = tagName
	}

	if options.PullImage {
		auth := findRegistry(options.Registries, imageSpec.registry)
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

	if options.DeleteContainer {
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
