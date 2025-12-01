package commands

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/podman"
	"github.com/danecwalker/otari/internal/spinners"
	"github.com/danecwalker/otari/internal/systemd"
	"github.com/danecwalker/otari/internal/utils"
	"github.com/fatih/color"
)

func Stop(ctx context.Context, stackPath string) {
	if stackPath == "" {
		stackPath = utils.DefaultStackPath()
	}
	c, err := os.ReadFile(stackPath)
	if err != nil {
		fmt.Println(utils.Error("Failed to read " + stackPath))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	stack, err := definition.Parse(c)
	if err != nil {
		fmt.Println(utils.Error("Failed to parse stack definition"))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	stack.StackName = utils.StackNameFromPath(stackPath)

	// Start all containers
	active, err := podman.ActiveContainers(ctx)
	if err != nil {
		fmt.Println(utils.Error("Failed to get active containers."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}
	for _, container := range stack.Containers {
		sp := spinners.DefaultSpinner()
		containerUnitName := container.ContainerName
		sp.SetMessage(fmt.Sprintf("Stopping container '%s'...", containerUnitName))

		// check if container is already running
		if isActive := slices.Contains(active, containerUnitName); !isActive {
			sp.FinishWithInfo(fmt.Sprintf("Container '%s' is already stopped.", containerUnitName))
			continue
		}

		if err := systemd.StopUnit(containerUnitName); err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to stop container '%s'", containerUnitName))
			// journalctl --user -xe -t portfolio
			// tell user to check journalctl for errors
			color.New(color.FgWhite).Println("    " + err.Error())
			color.New(color.FgWhite).Println("    Please check the container logs using 'journalctl --user -xe -t " + containerUnitName + "' for more details.")

			os.Exit(1)
		}
		sp.FinishWithSuccess(fmt.Sprintf("Container '%s' stopped.", containerUnitName))
	}
}
