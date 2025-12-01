package quadlets

import (
	"bytes"
	"fmt"

	"github.com/danecwalker/podstack/internal/definition"
	"github.com/danecwalker/podstack/internal/utils"
)

func (q *QuadletGenerator) GenerateVolume(stack *definition.Stack, volumeName string) ([]byte, error) {
	volume, exists := stack.Volumes[volumeName]
	if !exists {
		return nil, fmt.Errorf("volume '%s' not found in stack", volumeName)
	}

	var buf bytes.Buffer
	err := utils.WriteSection(&buf, "Unit", [][2]string{
		{"Description", fmt.Sprintf("%s volume", volume.VolumeName)},
	})
	if err != nil {
		return nil, err
	}

	if err := utils.WriteEmptyLine(&buf); err != nil {
		return nil, err
	}

	networkProperties := [][2]string{
		{"VolumeName", volume.VolumeName},
	}

	err = utils.WriteSection(&buf, "Volume", networkProperties)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
