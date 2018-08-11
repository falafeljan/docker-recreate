package recreate

import (
	docker "github.com/fsouza/go-dockerclient"
)

// RegistryConf describes authentication for private Docker registries.
type RegistryConf struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func findRegistry(
	registries []RegistryConf,
	registryHost string) (auth docker.AuthConfiguration) {
	for _, registry := range registries {
		if registry.Host == registryHost {
			return docker.AuthConfiguration{
				Username: registry.User,
				Password: registry.Password}
		}
	}

	return docker.AuthConfiguration{}
}
