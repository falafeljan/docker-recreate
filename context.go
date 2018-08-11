package recreate

import (
	"github.com/fsouza/go-dockerclient"
)

// DockerOptions describe additional options for pulling and creating the container
type DockerOptions struct {
	PullImage       bool
	DeleteContainer bool
	Registries      []RegistryConf
}

// Context describes the context needed for managing tokens
type Context struct {
	client  *docker.Client
	options DockerOptions
}

// NewContext creates a new recreation context with the default Docker environment
func NewContext(options DockerOptions) (context Context, err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return Context{}, err
	}

	return NewContextWithClient(options, client), nil
}

// NewContextWithEndpoint creates a new recreation context with a specified Docker endpoint
func NewContextWithEndpoint(options DockerOptions, endpoint string) (context Context, err error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return Context{}, err
	}

	return NewContextWithClient(options, client), nil
}

// NewContextWithClient creates a new recreation context with an existing Docker client
func NewContextWithClient(options DockerOptions, client *docker.Client) Context {
	return Context{
		client:  client,
		options: options,
	}
}
