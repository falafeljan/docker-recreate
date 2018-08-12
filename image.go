package recreate

import (
	"strings"
)

// ImageSpec describes all parts of an image identifier
type ImageSpec struct {
	registry   string
	name       string
	tag        string
	repository string
}

func generateImageURL(spec ImageSpec) string {
	if spec.tag == "" {
		return spec.repository
	}

	return spec.repository + ":" + spec.tag
}

func parseImageName(imageName string) (imageSpec ImageSpec) {
	registry := ""
	name := imageName
	tag := ""
	repository := ""

	slashIndex := strings.Index(name, "/")

	if slashIndex > -1 {
		registry = imageName[:slashIndex]
		name = imageName[(slashIndex + 1):]
	}

	colonIndex := strings.LastIndex(name, ":")

	if colonIndex > -1 {
		fullName := name
		name = fullName[:colonIndex]
		tag = fullName[(colonIndex + 1):]
	}

	if registry != "" {
		repository = strings.Join([]string{
			registry,
			name},
			"/")
	} else {
		repository = name
	}

	return ImageSpec{
		registry:   registry,
		name:       name,
		tag:        tag,
		repository: repository,
	}
}
