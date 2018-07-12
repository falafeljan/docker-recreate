package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	homedir "github.com/mitchellh/go-homedir"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func parseConf() (conf *Conf, err error) {
	emptyConf := Conf{Registries: []RegistryConf{}}
	homeDirectory, err := homedir.Dir()
	if err != nil {
		return &emptyConf, err
	}

	filePath := strings.Join([]string{
		homeDirectory,
		".recreate.json"},
		"/")

	file, err := os.Open(filePath)
	if err != nil {
		return &emptyConf, err
	}

	defer file.Close()

	var parsedConf Conf
	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, &parsedConf)

	return &parsedConf, nil
}

func createOptions(args *Args, conf *Conf) (options *Options) {
	return &Options{
		pullImage:       args.pullImage,
		deleteContainer: args.deleteContainer,
		registries:      conf.Registries}
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
	checkError(err)

	conf, _ := parseConf()
	checkError(err)

	recreation, err := RecreateWithClient(
		client,
		args.containerID,
		args.tagName,
		createOptions(&args, conf))
	checkError(err)

	fmt.Printf(
		"Migrated `%s` from %s to %s.\n",
		args.containerID,
		recreation.previousContainerID[:4],
		recreation.newContainerID[:4])
}
