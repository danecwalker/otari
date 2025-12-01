package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestContainerImageName(t *testing.T) {
	tests := []struct {
		imageName string
		valid     bool
	}{
		{"nginx:latest", false},
		{"redis", false},
		{"myregistry.local:5000/test/image:tag", true},
		{"invalid image name", false},
		{"another@invalid:image", false},
		{"myregistry.local/test/image", true},
		{"myregistry.local/test/image:latest", true},
		{"myregistry.local/test/image@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", true},
	}

	for _, tt := range tests {
		t.Run(tt.imageName, func(t *testing.T) {
			var im Image
			yaml.Unmarshal([]byte(tt.imageName), &im)
			assert.Equal(t, tt.valid, im.FullyQual)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		imageRef *Image
		expected string
	}{
		{
			imageRef: &Image{
				Registry:  "myregistry.local:5000",
				Project:   "test",
				Image:     "image",
				Tag:       "tag",
				FullyQual: true,
			},
			expected: "myregistry.local:5000/test/image:tag",
		},
		{
			imageRef: &Image{
				Registry:  "myregistry.local",
				Project:   "test",
				Image:     "image",
				Digest:    "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				FullyQual: true,
			},
			expected: "myregistry.local/test/image@sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			imageRef: &Image{
				Registry: "localhost",
				Image:    "nginx",
			},
			expected: "localhost/nginx",
		},
		{
			imageRef: &Image{
				Image: "nginx",
			},
			expected: "nginx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.imageRef.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}
