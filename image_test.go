package recreate

import "testing"

func TestParseImageName(t *testing.T) {
	fixtures := []string{
		"foobar",
		"foobar:2",
		"foobar:latest",
		"registry.acme.corp/foobar",
		"registry.acme.corp/foobar:2",
		"registry.acme.corp/foobar:latest"}

	expecteds := [][]string{
		[]string{"", "foobar", "latest"},
		[]string{"", "foobar", "2"},
		[]string{"", "foobar", "latest"},
		[]string{"registry.acme.corp", "foobar", "latest"},
		[]string{"registry.acme.corp", "foobar", "2"},
		[]string{"registry.acme.corp", "foobar", "latest"}}

	for i, fixture := range fixtures {
		expected := expecteds[i]
		imageSpec := parseImageName(fixture)

		if imageSpec.registry != expected[0] || imageSpec.name != expected[1] || imageSpec.tag != expected[2] {
			t.Errorf("Image structure does not equal:\nExpected: %v\nReceived: %v %v %v\n",
				expected,
				imageSpec.registry,
				imageSpec.name,
				imageSpec.tag)
		}
	}
}
