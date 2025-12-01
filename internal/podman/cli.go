package podman

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

func PodmanVersion() (bool, string) {
	cmd := exec.Command("podman", "version", "-f", "{{ .Version }}")
	out, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	return true, string(bytes.TrimSpace(out))
}

func ParsePodmanVersion(version string) (int, int, int) {
	var major, minor, patch int
	n, err := fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	if err != nil || n < 2 {
		return 0, 0, 0
	}
	return major, minor, patch
}

func ActiveContainers(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "podman", "ps", "--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(bytes.TrimSpace(out), []byte{'\n'})
	var containers []string
	for _, line := range lines {
		containers = append(containers, string(line))
	}
	return containers, nil
}

func RemoveNetwork(ctx context.Context, networkName string) error {
	cmd := exec.CommandContext(ctx, "podman", "network", "rm", "-f", networkName)
	return cmd.Run()
}

func RemoveVolume(ctx context.Context, volumeName string) error {
	cmd := exec.CommandContext(ctx, "podman", "volume", "rm", "-f", volumeName)
	return cmd.Run()
}
