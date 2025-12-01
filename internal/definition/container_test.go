package definition

// Tests the unmarshalling of Container structures from YAML.

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestContainerUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		expected *Container
	}{
		{
			name: "Basic container with image and environment as map",
			yamlData: `image: nginx:latest
environment:
  VAR1: value1	
  VAR2: value2`,
			expected: &Container{
				Image: &Image{
					Image:     "nginx",
					Tag:       "latest",
					FullyQual: false,
				},
				Environment: MapArray{
					"VAR1": "value1",
					"VAR2": "value2",
				},
			},
		},
		{
			name: "Container with environment as array",
			yamlData: `image: redis:alpine
environment:
  - VAR1=value1
  - VAR2=value2`,
			expected: &Container{
				Image: &Image{
					Image:     "redis",
					Tag:       "alpine",
					FullyQual: false,
				},
				Environment: MapArray{
					"VAR1": "value1",
					"VAR2": "value2",
				},
			},
		},
		{
			name:     "Container with no environment",
			yamlData: `image: postgres:latest`,
			expected: &Container{
				Image: &Image{
					Image:     "postgres",
					Tag:       "latest",
					FullyQual: false,
				},
				Environment: nil,
			},
		},
		{
			name: "Container with entrypoint as single string",
			yamlData: `image: alpine:latest
entrypoint: /bin/sh`,
			expected: &Container{
				Image: &Image{
					Image:     "alpine",
					Tag:       "latest",
					FullyQual: false,
				},
				Entrypoint: StringArray("/bin/sh"),
			},
		},
		{
			name: "Container with entrypoint as array",
			yamlData: `image: alpine:latest
entrypoint:
  - /bin/sh
  - -c`,
			expected: &Container{
				Image: &Image{
					Image:     "alpine",
					Tag:       "latest",
					FullyQual: false,
				},
				Entrypoint: StringArray("/bin/sh -c"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var container Container
			err := yaml.Unmarshal([]byte(tt.yamlData), &container)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, &container)
		})
	}
}
