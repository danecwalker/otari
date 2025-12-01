package rules

import (
	"fmt"
	"strings"

	"github.com/danecwalker/otari/internal/definition"
)

func ValidateContainerNames(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	nameSet := make(map[string]struct{})

	for _, container := range s.Containers {
		if _, exists := nameSet[container.ContainerName]; exists {
			errors = append(errors, &RuleError{
				Message: "Duplicate container name '" + container.ContainerName + "' found.",
			})
		} else {
			nameSet[container.ContainerName] = struct{}{}
		}
	}

	return errors
}

func ValidateDuplicateEnvironmentVariables(s *definition.Stack) []*RuleError {
	var errors []*RuleError

	for _, container := range s.Containers {
		envSet := make(map[string]struct{})
		for envKey := range container.Environment {
			if _, exists := envSet[envKey]; exists {
				errors = append(errors, &RuleError{
					Message: "Duplicate environment variable '" + envKey + "' found in container '" + container.ContainerName + "'.",
				})
			} else {
				envSet[envKey] = struct{}{}
			}
		}
	}

	return errors
}

func ValidateDependencyExistence(s *definition.Stack) []*RuleError {
	var errors []*RuleError
	for _, container := range s.Containers {
		for _, depName := range container.Depends {
			if _, exists := s.Containers[depName]; !exists {
				errors = append(errors, &RuleError{
					Message: "Container '" + container.ContainerName + "' has undefined dependency '" + depName + "'.",
				})
			}
		}
	}
	return errors
}

func ValidateCircularDependencies(s *definition.Stack) []*RuleError {
	var errors []*RuleError

	// Build adjacency: containerName -> []dependsOnNames
	deps := make(map[string][]string, len(s.Containers))
	for cname, c := range s.Containers {
		deps[cname] = append([]string(nil), c.Depends...)
	}

	// 0 = unvisited, 1 = visiting, 2 = done
	const (
		stateUnvisited = 0
		stateVisiting  = 1
		stateDone      = 2
	)

	state := make(map[string]int, len(deps))
	seenCycles := make(map[string]bool) // to avoid duplicate reports

	var dfs func(name string, stack []string)

	dfs = func(name string, stack []string) {
		switch state[name] {
		case stateDone:
			return
		case stateVisiting:
			// Found a back edge → cycle.
			// Extract the cycle from the stack.
			startIdx := 0
			for i, n := range stack {
				if n == name {
					startIdx = i
					break
				}
			}
			cycle := append(stack[startIdx:], name) // close the loop visually
			key := strings.Join(cycle, "->")
			if !seenCycles[key] {
				seenCycles[key] = true
				errors = append(errors, &RuleError{
					Message: fmt.Sprintf("circular dependency detected: %s",
						strings.Join(cycle, " -> ")),
				})
			}
			return
		}

		state[name] = stateVisiting
		stack = append(stack, name)

		for _, dep := range deps[name] {
			// If dependency refers to an unknown container, ignore here;
			// that's the job of another rule (ValidateContainerDependenciesExistence, etc).
			if _, ok := deps[dep]; !ok {
				continue
			}
			dfs(dep, stack)
		}

		state[name] = stateDone
	}

	// Run DFS from each container
	for name := range deps {
		if state[name] == stateUnvisited {
			dfs(name, nil)
		}
	}

	return errors
}
