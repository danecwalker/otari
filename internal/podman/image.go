package podman

import (
	"context"
	"os/exec"
)

func ImageExists(ctx context.Context, image string) bool {
	_, err := exec.CommandContext(ctx, "podman", "image", "exists", image).Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			status := exitError.ExitCode()
			if status == 1 {
				return false
			}
			return false
		}
		return false
	}
	return true
}

func ImagePull(ctx context.Context, image string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "podman", "pull", image)
	return cmd
}
