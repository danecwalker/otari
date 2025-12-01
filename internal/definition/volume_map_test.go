package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestVolumeMap(t *testing.T) {
	tests := []struct {
		vMYaml   string
		expected string
	}{
		{
			vMYaml: `
volumes:
      - ./:/code`,
			expected: "./:/code",
		},
	}

	for _, tt := range tests {
		var d struct {
			Volumes []VolumeMap `yaml:"volumes"`
		}
		err := yaml.Unmarshal([]byte(tt.vMYaml), &d)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, d.Volumes[0].String())
	}
}
