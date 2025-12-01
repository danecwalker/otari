package quadlets

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/danecwalker/podstack/internal/definition"
	"github.com/danecwalker/podstack/internal/utils"
)

// GenerateContainer implements generate.Generator.
func (q *QuadletGenerator) GenerateContainer(stack *definition.Stack, containerName string) ([]byte, error) {
	container, exists := stack.Containers[containerName]
	if !exists {
		return nil, fmt.Errorf("container '%s' not found in stack", containerName)
	}
	var buf bytes.Buffer

	unitDefinition := [][2]string{
		{"Description", container.ContainerName + " container"},
	}

	if len(container.Depends) > 0 {
		unitDefinition = append(unitDefinition, [2]string{
			"Requires", strings.Join(container.Depends, " "),
		})
		unitDefinition = append(unitDefinition, [2]string{
			"After", strings.Join(container.Depends, " "),
		})
	}

	if err := utils.WriteSection(&buf, "Unit", unitDefinition); err != nil {
		return nil, err
	}

	if err := utils.WriteEmptyLine(&buf); err != nil {
		return nil, err
	}

	// Container section
	containerProperties := [][2]string{
		{"ContainerName", container.ContainerName},
		{"Image", container.Image.String()},
	}

	if container.Entrypoint != "" {
		containerProperties = append(containerProperties, [2]string{
			"EntryPoint", string(container.Entrypoint),
		})
	}

	if container.Init {
		containerProperties = append(containerProperties, [2]string{
			"RunInit", "true",
		})
	}

	// Environment variables
	if len(container.Environment) > 0 {
		for key, value := range container.Environment {
			containerProperties = append(containerProperties, [2]string{
				"Environment", key + "=" + value,
			})
		}

	}

	if len(container.Ports) > 0 {
		var portMappings []string
		for _, port := range container.Ports {
			portMappings = append(portMappings, port.String())
		}
		containerProperties = append(containerProperties, [2]string{
			"PublishPort", strings.Join(portMappings, " "),
		})
	}

	if len(container.Labels) > 0 {
		for key, value := range container.Labels {
			containerProperties = append(containerProperties, [2]string{
				"Label", key + "=" + value,
			})
		}
	}

	// check if contains any host networks
	if len(container.Networks) > 0 {
		hasHostNetwork := false
		for _, network := range container.Networks {
			if stack.Networks[network].Driver == definition.NetworkDriverHost {
				hasHostNetwork = true
				break
			}
		}

		if hasHostNetwork {
			containerProperties = append(containerProperties, [2]string{
				"Network", "host",
			})
		} else {
			for _, network := range container.Networks {
				containerProperties = append(containerProperties, [2]string{
					"Network", network + ".network",
				})
			}
		}
	}

	if len(container.Volumes) > 0 {
		for _, volumeMap := range container.Volumes {
			volumeDef := volumeMap.Destination + ":" + strings.Join(volumeMap.Options, ",")
			if volumeMap.Type == definition.VolumeMountTypeBind {
				// get absolute path for host bind mounts
				absPath, err := utils.GetAbsolutePath(volumeMap.Source)
				if err != nil {
					return nil, fmt.Errorf("failed to get absolute path for bind mount '%s': %v", volumeMap.Source, err)
				}
				volumeDef = absPath + ":" + volumeDef
			} else {
				volumeDef = volumeMap.Source + ".volume:" + volumeDef
			}
			containerProperties = append(containerProperties, [2]string{
				"Volume", volumeDef,
			})
		}
	}

	err := utils.WriteSection(&buf, "Container", containerProperties)
	if err != nil {
		return nil, err
	}

	if err := utils.WriteEmptyLine(&buf); err != nil {
		return nil, err
	}

	serviceProperties := [][2]string{
		{"TimeoutStartSec", "900"}, // Simplified for example purposes
	}

	// Service section
	if container.Restart.IsNo() {
		serviceProperties = append(serviceProperties, [2]string{
			"Restart", "no",
		})
	} else if container.Restart.IsAlways() {
		serviceProperties = append(serviceProperties, [2]string{
			"Restart", "always",
		})
	} else if container.Restart.IsOnFailure() {
		serviceProperties = append(serviceProperties, [2]string{
			"Restart", "on-failure",
		})
	} else if container.Restart.IsUnlessStopped() {
		serviceProperties = append(serviceProperties, [2]string{
			"Restart", "always",
		})

		serviceProperties = append(serviceProperties, [2]string{
			"RestartPreventExitStatus", "0 SIGKILL",
		})
	}

	err = utils.WriteSection(&buf, "Service", serviceProperties)
	if err != nil {
		return nil, err
	}

	if err := utils.WriteEmptyLine(&buf); err != nil {
		return nil, err
	}
	// Install section
	installProperties := [][2]string{
		{"WantedBy", "multi-user.target default.target"},
	}

	err = utils.WriteSection(&buf, "Install", installProperties)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
