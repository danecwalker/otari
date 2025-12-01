package podman

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/utils"
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

func ImageBuild(ctx context.Context, build *definition.Build, image string) (*exec.Cmd, error) {
	// check if build context path exists
	absPath, err := utils.GetAbsolutePath(build.Context)
	if err != nil {
		return nil, err
	}

	if !utils.PathExists(absPath) {
		return nil, fmt.Errorf("build context path does not exist: %s", absPath)
	}

	// check if containerfile / dockerfile exists
	containerFile := build.ContainerFile
	if containerFile == "" {
		containerFile = "Containerfile"

		if !utils.PathExists(filepath.Join(absPath, containerFile)) {
			containerFile = "Dockerfile"
		}
	}

	// check custom containerfile path
	if !filepath.IsAbs(containerFile) {
		containerFile = filepath.Join(absPath, containerFile)
	}
	if !utils.PathExists(containerFile) {
		return nil, fmt.Errorf("containerfile does not exist: %s", containerFile)
	}

	cmdSlice := []string{
		"build",
		"-t", image,
		"-f", containerFile,
	}

	// build args
	for key, value := range build.Args {
		cmdSlice = append(cmdSlice, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	}

	// build tags
	for _, tag := range build.Tags {
		cmdSlice = append(cmdSlice, "-t", tag)
	}

	if build.Target != "" {
		cmdSlice = append(cmdSlice, "--target", build.Target)
	}

	cmd := exec.CommandContext(ctx, "podman", append(cmdSlice, absPath)...)
	return cmd, nil
}
