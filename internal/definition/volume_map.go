package definition

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var volumeRe = regexp.MustCompile(`^([^:]+):([^:]+)(?::((?:rw|ro|z|Z)(?:,(?:rw|ro|z|Z))*))?$`)

type VolumeMountType string

const (
	VolumeMountTypeBind   VolumeMountType = "bind"
	VolumeMountTypeVolume VolumeMountType = "volume"
)

type VolumeMap struct {
	Source      string
	Destination string
	Options     []string
	Type        VolumeMountType
}

func ParseVolumeMap(volumeStr string) (*VolumeMap, error) {
	matches := volumeRe.FindStringSubmatch(volumeStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid volume format: %s", volumeStr)
	}

	source := matches[1]
	destination := matches[2]
	var options []string
	if matches[3] != "" {
		options = strings.Split(matches[3], ",")
	}

	return &VolumeMap{
		Source:      source,
		Destination: destination,
		Options:     options,
		Type:        VolumeMountTypeVolume,
	}, nil
}

func (vm *VolumeMap) String() string {
	result := fmt.Sprintf("%s:%s", vm.Source, vm.Destination)
	if len(vm.Options) > 0 {
		result += ":" + strings.Join(vm.Options, ",")
	}
	return result
}

func (vm *VolumeMap) UnmarshalYAML(value *yaml.Node) error {
	var volumeStr string
	if err := value.Decode(&volumeStr); err != nil {
		return err
	}

	parsedVM, err := ParseVolumeMap(volumeStr)
	if err != nil {
		return err
	}

	*vm = *parsedVM
	return nil
}
