package rules

import (
	"fmt"

	"github.com/danecwalker/podstack/internal/definition"
)

func ValidateNetworkNames(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	nameSet := make(map[string]struct{})

	for _, network := range s.Networks {
		if _, exists := nameSet[network.NetworkName]; exists {
			errors = append(errors, &RuleError{
				Message: "Duplicate network name '" + network.NetworkName + "' found.",
			})
		} else {
			nameSet[network.NetworkName] = struct{}{}
		}
	}

	return errors
}

func ValidateContainerNetworkExistence(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	for _, container := range s.Containers {
		for _, networkName := range container.Networks {
			if _, exists := s.Networks[networkName]; !exists {
				errors = append(errors, &RuleError{
					Message: "Container '" + container.ContainerName + "' references undefined network '" + networkName + "'.",
				})
			}
		}
	}
	return errors
}

func ValidatePortConflicts(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	portSet := make(map[int]string)

	for _, container := range s.Containers {
		for _, port := range container.Ports {
			hostPort := port.HostPort
			if hostPort.Range {
				for p := hostPort.Start; p <= hostPort.End; p++ {
					if existingContainer, exists := portSet[p]; exists {
						errors = append(errors, &RuleError{
							Message: "Port conflict on port " + fmt.Sprint(p) + " between containers '" + existingContainer + "' and '" + container.ContainerName + "'.",
						})
					} else {
						portSet[p] = container.ContainerName
					}
				}
			} else {
				p := hostPort.Start
				if existingContainer, exists := portSet[p]; exists {
					errors = append(errors, &RuleError{
						Message: "Port conflict on port " + fmt.Sprint(p) + " between containers '" + existingContainer + "' and '" + container.ContainerName + "'.",
					})
				} else {
					portSet[p] = container.ContainerName
				}
			}
		}
	}

	return errors
}

func ValidateHostNetworkPortConflicts(s *definition.Stack) []*RuleError {
	var errors []*RuleError

	for _, container := range s.Containers {
		// find if container uses host network
		for _, networkName := range container.Networks {
			network, exists := s.Networks[networkName]
			if exists && network.Driver == "host" {
				if len(container.Ports) > 0 {
					errors = append(errors, &RuleError{
						Message: "Container '" + container.ContainerName + "' uses host network and defines port mappings, which is a conflict.",
					})
				}
			}
		}
	}

	return errors
}
