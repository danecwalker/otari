package podman

import (
	"bytes"
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
