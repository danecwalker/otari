package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/danecwalker/podstack/internal/definition"
)

func ValidateVolumeNames(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	nameSet := make(map[string]struct{})

	for _, volume := range s.Volumes {
		if _, exists := nameSet[volume.VolumeName]; exists {
			errors = append(errors, &RuleError{
				Message: "Duplicate volume name '" + volume.VolumeName + "' found.",
			})
		} else {
			nameSet[volume.VolumeName] = struct{}{}
		}
	}

	return errors
}

func isHostPath(p string) bool {
	// Absolute path (POSIX or Windows)
	if filepath.IsAbs(p) {
		return true
	}

	// Contains a path separator → Docker says it's a host path
	if strings.ContainsRune(p, os.PathSeparator) {
		return true
	}

	// Windows alt separator `/` still counts as host path
	if os.PathSeparator != '/' && strings.ContainsRune(p, '/') {
		return true
	}

	return false
}

func ValidateContainerVolumeExistence(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	for cname, container := range s.Containers {
		for i, volumeMap := range container.Volumes {
			volumeName := volumeMap.Source
			if _, exists := s.Volumes[volumeName]; !exists {
				// Also check for the case where the volume is specified as a host path
				if isHostPath(volumeName) {
					// check if the path exists on the host
					if _, err := os.Stat(volumeName); err == nil {
						// change volume mount type to bind mount
						s.Containers[cname].Volumes[i].Type = definition.VolumeMountTypeBind
						continue
					}
				}
				errors = append(errors, &RuleError{
					Message: "Container '" + container.ContainerName + "' references undefined volume '" + volumeName + "'.",
				})
			}
		}
	}
	return errors
}

func ValidateDuplicateVolumeMountsPerContainer(s *definition.Stack) []*RuleError {
	var errors []*RuleError

	for _, container := range s.Containers {
		mountSet := make(map[string]struct{})
		for _, volumeMap := range container.Volumes {
			mountPoint := volumeMap.Destination
			if _, exists := mountSet[mountPoint]; exists {
				errors = append(errors, &RuleError{
					Message: "Container '" + container.ContainerName + "' has duplicate volume mount point '" + mountPoint + "'.",
				})
			} else {
				mountSet[mountPoint] = struct{}{}
			}
		}
	}

	return errors
}
