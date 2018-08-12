package recreate

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"math/rand"
	"testing"
	"time"
)

var imageName = "recreate_testing"
var dockerfile = []byte(
	"FROM scratch\n" +
		"CMD tail -f /dev/null\n")

// <https://stackoverflow.com/a/22892986>
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func buildImage(client *docker.Client, imageName string, imageTags []string) (err error) {
	t := time.Now()
	inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	tr := tar.NewWriter(inputbuf)
	tr.WriteHeader(&tar.Header{
		Name:       "Dockerfile",
		Size:       int64(len(dockerfile)),
		ModTime:    t,
		AccessTime: t,
		ChangeTime: t,
	})
	tr.Write(dockerfile)
	tr.Close()

	opts := docker.BuildImageOptions{
		Name:         imageName,
		InputStream:  inputbuf,
		OutputStream: outputbuf,
	}

	if err := client.BuildImage(opts); err != nil {
		return err
	}

	for _, imageTag := range imageTags {
		client.TagImage(imageName, docker.TagImageOptions{
			Repo: imageName,
			Tag:  imageTag,
		})
	}

	return nil
}

func performRecreationWithOptions(
	client *docker.Client,
	imageTag string,
	dockerOptions DockerOptions,
	containerOptions ContainerOptions,
	f func(*Recreation) error,
) error {
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: "recreate_testing_" + randSeq(16),
		Config: &docker.Config{
			Image: imageName,
		},
	})

	if err != nil {
		return err
	}
	fmt.Printf("Container: %s", container.ID)
	err = client.StartContainer(container.ID, container.HostConfig)
	if err != nil {
		return err
	}

	context := NewContextWithClient(dockerOptions, client)

	res, err := context.Recreate(
		container.ID,
		imageTag,
		containerOptions,
	)

	if err != nil {
		return err
	}

	err = f(res)
	if err != nil {
		return err
	}

	return nil
}

func TestRecreate(t *testing.T) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	err = buildImage(client, imageName, []string{"foo", "bar"})
	if err != nil {
		t.Fatal(err)
	}

	// Test whether it is recreated when no tag is specified
	err = performRecreationWithOptions(
		client,
		"",
		DockerOptions{DeleteContainer: true},
		ContainerOptions{},
		func(res *Recreation) error {
			container, err := client.InspectContainer(res.NewContainerID)
			if err != nil {
				return err
			}

			if container.Image != imageName {
				return fmt.Errorf("Container image does not match: `%s` vs `%s`", imageName, container.Image)
			}

			return nil
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	// Test whether old container stays if told so
	// t.Fatal(errors.New("Not implemented"))

	// Test whether old container gets deleted if told so
	// t.Fatal(errors.New("Not implemented"))

	// Test whether environment variables are applied to new container
	// t.Fatal(errors.New("Not implemented"))
}
