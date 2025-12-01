package definition

import (
	"gopkg.in/yaml.v3"
)

type Stack struct {
	Containers map[string]*Container `yaml:"containers"`
	Volumes    map[string]*Volume    `yaml:"volumes"`
	Networks   map[string]*Network   `yaml:"networks"`
}

func Parse(data []byte) (*Stack, error) {
	var s Stack
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	for name, container := range s.Containers {
		container.ContainerName = name
	}

	for name, volume := range s.Volumes {
		if volume == nil {
			volume = &Volume{}
			s.Volumes[name] = volume
		}

		volume.VolumeName = name
	}

	for name, network := range s.Networks {
		if network == nil {
			network = &Network{}
			s.Networks[name] = network
		}
		network.NetworkName = name
	}

	return &s, nil
}
