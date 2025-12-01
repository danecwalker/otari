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

func Remove(ctx context.Context, stackPath string) {
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

	// Stop all containers
	active, err := podman.ActiveContainers(ctx)
	if err != nil {
		fmt.Println(utils.Error("Failed to get active containers."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}
	for _, container := range stack.Containers {
		sp := spinners.DefaultSpinner()
		containerUnitName := container.ContainerName
		sp.SetMessage(fmt.Sprintf("Removing container '%s'...", containerUnitName))

		// check if container is already running
		if isActive := slices.Contains(active, containerUnitName); !isActive {
			sp.Println(color.New(color.FgWhite).Sprintf("Container '%s' is already stopped.", containerUnitName))
			if err := systemd.DeleteUnitFile(containerUnitName + ".container"); err != nil {
				sp.FinishWithError(fmt.Sprintf("Failed to remove container unit file for '%s'", containerUnitName))
				color.New(color.FgWhite).Println("    " + err.Error())
				os.Exit(1)
			}
			sp.FinishWithSuccess(fmt.Sprintf("Container '%s' removed.", containerUnitName))
			continue
		}

		sp.Println(color.New(color.FgWhite).Sprintf("Container '%s' is running, stopping it first.", containerUnitName))
		if err := systemd.StopUnit(containerUnitName); err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to stop container '%s'", containerUnitName))
			// journalctl --user -xe -t portfolio
			// tell user to check journalctl for errors
			color.New(color.FgWhite).Println("    " + err.Error())
			color.New(color.FgWhite).Println("    Please check the container logs using 'journalctl --user -xe -t " + containerUnitName + "' for more details.")

			os.Exit(1)
		}

		sp.Println(color.New(color.FgWhite).Sprintf("Container '%s' stopped, removing unit file.", containerUnitName))

		if err := systemd.DeleteUnitFile(containerUnitName + ".container"); err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to remove container unit file for '%s'", containerUnitName))
			color.New(color.FgWhite).Println("    " + err.Error())
			os.Exit(1)
		}
		sp.FinishWithSuccess(fmt.Sprintf("Container '%s' removed.", containerUnitName))
	}

	active, err = podman.ActiveContainers(ctx)
	if err != nil {
		fmt.Println(utils.Error("Failed to get active containers."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}

	// Remove volumes
	for _, volume := range stack.Volumes {
		if !volume.PersistOnRemove {
			sp := spinners.DefaultSpinner()
			volumeUnitName := volume.VolumeName
			sp.SetMessage(fmt.Sprintf("Removing volume '%s'...", volumeUnitName))
			// Check if volume is in use by any active container
			volumeUsed := false
			for _, container := range stack.Containers {
				for _, vol := range container.Volumes {
					if vol.Source == volume.VolumeName && slices.Contains(active, container.ContainerName) {
						volumeUsed = true
						break
					}
				}
			}
			if volumeUsed {
				sp.FinishWithInfo(fmt.Sprintf("Volume '%s' is still in use.", volumeUnitName))
			} else {
				// Remove the volume quadlet
				sp.SetMessage(fmt.Sprintf("Removing volume '%s'...", volumeUnitName))

				if err := systemd.StopUnit(volumeUnitName + "-volume"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to stop volume '%s'.", volumeUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				if err := systemd.DeleteUnitFile(volumeUnitName + ".volume"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove volume '%s'.", volumeUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				if err := podman.RemoveVolume(ctx, volumeUnitName); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove volume '%s' from Podman.", volumeUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				sp.FinishWithSuccess(fmt.Sprintf("Volume '%s' removed.", volumeUnitName))
			}
		}
	}

	// Remove networks
	for _, network := range stack.Networks {
		if !network.PersistOnRemove {
			sp := spinners.DefaultSpinner()
			networkUnitName := network.NetworkName
			sp.SetMessage(fmt.Sprintf("Removing network '%s'...", networkUnitName))
			// Check if network is in use by any active container
			networkUsed := false
			for _, container := range stack.Containers {
				for _, net := range container.Networks {
					if net == network.NetworkName && slices.Contains(active, container.ContainerName) {
						networkUsed = true
						break
					}
				}
			}
			if networkUsed {
				sp.FinishWithInfo(fmt.Sprintf("Network '%s' is still in use.", networkUnitName))
			} else {
				// Remove the network quadlet
				sp.SetMessage(fmt.Sprintf("Removing network '%s'...", networkUnitName))

				if err := systemd.StopUnit(networkUnitName + "-network"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to stop network '%s'.", networkUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				if err := systemd.DeleteUnitFile(networkUnitName + ".network"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove network '%s'.", networkUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				if err := podman.RemoveNetwork(ctx, networkUnitName); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove network '%s' from Podman.", networkUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}

				sp.FinishWithSuccess(fmt.Sprintf("Network '%s' removed.", networkUnitName))
			}
		}
	}

	// reload systemd daemon to apply changes
	if err := systemd.ReloadDaemon(); err != nil {
		fmt.Println(utils.Error("Failed to reload systemd daemon."))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	// remove lock file
	lockPath := stack.StackName + ".lock"
	if err := os.Remove(lockPath); err != nil {
		// if the file does not exist, ignore the error
		if !os.IsNotExist(err) {
			fmt.Println(utils.Error("Failed to remove stack lock file."))
			color.New(color.FgWhite).Println("    " + err.Error())
			os.Exit(1)
		}
	}

	fmt.Println(utils.Success("Stack '" + stack.StackName + "' removed successfully."))
}
