package recreate

import "testing"

func TestGenerateImageURL(t *testing.T) {
	fixtures := []ImageSpec{
		ImageSpec{
			tag:        "",
			repository: "registry.acme.corp/foobar",
		},
		ImageSpec{
			tag:        "latest",
			repository: "registry.acme.corp/foobar",
		},
	}

	expecteds := []string{
		"registry.acme.corp/foobar",
		"registry.acme.corp/foobar:latest",
	}

	for i, fixture := range fixtures {
		expected := expecteds[i]
		received := generateImageURL(fixture)

		if expected != received {
			t.Errorf("Image URL does not equal:\nExpected: %s\nReceived: %s\n", expected, received)
		}
	}
}

func TestParseImageName(t *testing.T) {
	fixtures := []string{
		"foobar",
		"foobar:2",
		"foobar:latest",
		"registry.acme.corp/foobar",
		"registry.acme.corp/foobar:2",
		"registry.acme.corp/foobar:latest",
	}

	expecteds := [][]string{
		[]string{"", "foobar", "", "foobar"},
		[]string{"", "foobar", "2", "foobar"},
		[]string{"", "foobar", "latest", "foobar"},
		[]string{"registry.acme.corp", "foobar", "", "registry.acme.corp/foobar"},
		[]string{"registry.acme.corp", "foobar", "2", "registry.acme.corp/foobar"},
		[]string{"registry.acme.corp", "foobar", "latest", "registry.acme.corp/foobar"},
	}

	for i, fixture := range fixtures {
		expected := expecteds[i]
		imageSpec := parseImageName(fixture)

		if imageSpec.registry != expected[0] ||
			imageSpec.name != expected[1] ||
			imageSpec.tag != expected[2] ||
			imageSpec.repository != expected[3] {
			t.Errorf("Image structure does not equal:\nExpected: %v\nReceived: %v %v %v %v\n",
				expected,
				imageSpec.registry,
				imageSpec.name,
				imageSpec.tag,
				imageSpec.repository,
			)
		}
	}
}
