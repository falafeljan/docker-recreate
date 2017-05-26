package main

import (
  "strings"
)

func parseImageName(imageName string) (repository string, tag string) {
  sepIndex := strings.LastIndex(imageName, ":")

  if sepIndex > -1 {
    repository := imageName[:sepIndex]
    tag := imageName[(sepIndex+1):]

    return repository, tag
  }

  return imageName, "latest"
}
