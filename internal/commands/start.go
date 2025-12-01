package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/danecwalker/otari/internal/changes"
	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/generate"
	"github.com/danecwalker/otari/internal/podman"
	"github.com/danecwalker/otari/internal/quadlets"
	"github.com/danecwalker/otari/internal/rules"
	"github.com/danecwalker/otari/internal/spinners"
	"github.com/danecwalker/otari/internal/systemd"
	"github.com/danecwalker/otari/internal/utils"
	"github.com/fatih/color"
)

func Start(ctx context.Context, stackPath string) {
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

	errors := rules.Validate(stack)
	if len(errors) > 0 {
		fmt.Println(utils.Error("Failed to validate stack:"))
		for _, err := range errors {
			color.New(color.FgRed).Printf("    • %s\n", err.Message)
		}
		os.Exit(1)
	}

	fmt.Println(utils.Success("Stack validated successfully!"))

	sp := spinners.DefaultSpinner()
	sp.SetMessage("Detecting changes...")
	new, deleted, totalChanges, err := changes.DetectChanges(ctx, stack)
	if err != nil {
		sp.FinishWithError("Failed to detect changes.")
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}
	if totalChanges > 0 {
		sp.FinishWithInfo(fmt.Sprintf("Detected %d change(s).", totalChanges))
	} else if totalChanges == -1 {
		sp.FinishWithInfo("No existing stack found.")
	} else {
		sp.FinishWithSuccess("No changes detected.")
	}

	if totalChanges != 0 {
		// check if containers
		if len(stack.Containers) == 0 {
			fmt.Println(utils.Info("No containers defined in the stack."))
			return
		}

		imagesSet := make(map[string]struct {
			Remote bool
			Build  *definition.Build
		})
		for _, container := range new.Containers {
			if container.Build != nil {
				imagesSet[container.ContainerName] = struct {
					Remote bool
					Build  *definition.Build
				}{Remote: false, Build: container.Build}
			} else if container.Image != nil {
				imagesSet[container.Image.String()] = struct {
					Remote bool
					Build  *definition.Build
				}{Remote: true, Build: nil}
			}
		}

		for image := range imagesSet {
			sp := spinners.DefaultSpinner()

			// Check if image exists
			if !podman.ImageExists(ctx, image) {
				if imagesSet[image].Remote {
					sp.SetMessage(fmt.Sprintf("Pulling image '%s'", image))
					cmd := podman.ImagePull(ctx, image)
					hasStderr := true
					stderr, err := cmd.StderrPipe()
					if err != nil {
						hasStderr = false
					}

					if err := cmd.Start(); err != nil {
						sp.FinishWithError(fmt.Sprintf("Failed to start pulling image '%s'", image))
						color.New(color.FgWhite).Println("    " + err.Error())
						os.Exit(1)
					}

					// Read stderr for progress
					if hasStderr {
						scanner := bufio.NewScanner(stderr)
						for scanner.Scan() {
							line := scanner.Text()
							sp.Println(color.New(color.FgWhite).Sprintf(" >  %s", line))
						}
					}

					if err := cmd.Wait(); err != nil {
						sp.FinishWithError(fmt.Sprintf("Failed to pull image '%s'", image))
						color.New(color.FgWhite).Println("    " + err.Error())
						os.Exit(1)
					}

					sp.FinishWithSuccess(fmt.Sprintf("Pulled image '%s'.", image))
				} else {
					sp.SetMessage(fmt.Sprintf("Building image '%s'", image))
					cmd, err := podman.ImageBuild(ctx, imagesSet[image].Build, fmt.Sprintf("%s_%s", stack.StackName, image))
					if err != nil {
						sp.FinishWithError(fmt.Sprintf("Failed to prepare build for image '%s'", image))
						color.New(color.FgWhite).Println("    " + err.Error())
						os.Exit(1)
					}

					hasStderr := true
					stderr, err := cmd.StderrPipe()
					if err != nil {
						hasStderr = false
					}

					if err := cmd.Start(); err != nil {
						sp.FinishWithError(fmt.Sprintf("Failed to start building image '%s'", image))
						color.New(color.FgWhite).Println("    " + err.Error())
						os.Exit(1)
					}

					// Read stderr for progress
					if hasStderr {
						scanner := bufio.NewScanner(stderr)
						for scanner.Scan() {
							line := scanner.Text()
							sp.Println(color.New(color.FgWhite).Sprintf(" >  %s", line))
						}
					}

					if err := cmd.Wait(); err != nil {
						sp.FinishWithError(fmt.Sprintf("Failed to build image '%s'", image))
						color.New(color.FgWhite).Println("    " + err.Error())
						os.Exit(1)
					}

					sp.FinishWithSuccess(fmt.Sprintf("Built image '%s'.", image))
				}
			} else {
				sp.FinishWithInfo(fmt.Sprintf("Image '%s' already exists.", image))
			}
		}

		// Generate systemd quadlets
		if len(new.Containers)+len(new.Volumes)+len(new.Networks) == 0 {
			fmt.Println(utils.Info("No changes detected that require quadlet generation."))
		} else {
			fmt.Println(utils.Info("Generating systemd quadlets..."))
			outputDir := utils.OutputLocation()
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Println(utils.Error("Failed to create output directory."))
				color.New(color.FgWhite).Println("    " + err.Error())
				os.Exit(1)
			}

			if err := generate.Generate(stack, new, outputDir, quadlets.Generator()); err != nil {
				fmt.Println(utils.Error("Failed to generate systemd quadlets."))
				color.New(color.FgWhite).Println("    " + err.Error())

				os.Exit(1)
			}
		}

		// Stop and delete removed containers
		if deleted != nil {
			for _, container := range deleted.Containers {
				containerUnitName := container.ContainerName
				sp = spinners.DefaultSpinner()
				sp.SetMessage(fmt.Sprintf("Removing container '%s'...", containerUnitName))
				if err := container.Remove(); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove container '%s'.", containerUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}
				sp.FinishWithSuccess(fmt.Sprintf("Container '%s' removed.", containerUnitName))
			}

			// Check if container networks are used by other containers
			for _, network := range deleted.Networks {
				sp = spinners.DefaultSpinner()
				networkUsed := false
				for _, container := range stack.Containers {
					for _, net := range container.Networks {
						if net == network.NetworkName && deleted.Containers[container.ContainerName] == nil {
							networkUsed = true
							break
						}
					}
					if networkUsed {
						break
					}
				}
				if networkUsed {
					sp.SetMessage(fmt.Sprintf("Skipping removal of network '%s' as it is still in use.", network.NetworkName))
					sp.FinishWithInfo(fmt.Sprintf("Network '%s' is still in use.", network.NetworkName))
					continue
				}

				// Remove the network quadlet
				networkUnitName := network.NetworkName
				sp.SetMessage(fmt.Sprintf("Removing network '%s'...", networkUnitName))
				if err := systemd.DeleteUnitFile(networkUnitName + ".network"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove network '%s'.", networkUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}
				sp.FinishWithSuccess(fmt.Sprintf("Network '%s' removed.", networkUnitName))
			}

			// Check if container networks are used by other containers
			for _, volume := range deleted.Volumes {
				sp = spinners.DefaultSpinner()
				volumeUsed := false
				for _, container := range stack.Containers {
					for _, vol := range container.Volumes {
						if vol.Source == volume.VolumeName && deleted.Containers[container.ContainerName] == nil {
							volumeUsed = true
							break
						}
					}
					if volumeUsed {
						break
					}
				}
				if volumeUsed {
					sp.SetMessage(fmt.Sprintf("Skipping removal of volume '%s' as it is still in use.", volume.VolumeName))
					sp.FinishWithInfo(fmt.Sprintf("Volume '%s' is still in use.", volume.VolumeName))
					continue
				}

				// Remove the volume quadlet
				volumeUnitName := volume.VolumeName
				sp.SetMessage(fmt.Sprintf("Removing volume '%s'...", volumeUnitName))
				if err := systemd.DeleteUnitFile(volumeUnitName + ".volume"); err != nil {
					sp.FinishWithError(fmt.Sprintf("Failed to remove volume '%s'.", volumeUnitName))
					color.New(color.FgWhite).Println("    " + err.Error())
					os.Exit(1)
				}
				sp.FinishWithSuccess(fmt.Sprintf("Volume '%s' removed.", volumeUnitName))
			}
		}

	}

	fmt.Println()

	if err := systemd.ReloadDaemon(); err != nil {
		fmt.Println(utils.Error("Failed to reload systemd daemon."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}

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
		sp.SetMessage(fmt.Sprintf("Starting container '%s'...", containerUnitName))

		// check if container is already running
		if isActive := slices.Contains(active, containerUnitName); isActive {
			sp.FinishWithInfo(fmt.Sprintf("Container '%s' is already running.", containerUnitName))
			continue
		}

		if err := systemd.StartUnit(containerUnitName); err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to start container '%s'", containerUnitName))
			// journalctl --user -xe -t portfolio
			// tell user to check journalctl for errors
			color.New(color.FgWhite).Println("    " + err.Error())
			color.New(color.FgWhite).Println("    Please check the container logs using 'journalctl --user -xe -t " + containerUnitName + "' for more details.")

			os.Exit(1)
		}
		sp.FinishWithSuccess(fmt.Sprintf("Container '%s' started.", containerUnitName))
	}

	sp = spinners.DefaultSpinner()
	sp.SetMessage("Computing change hashes...")
	if err := changes.SaveStackData(stack); err != nil {
		sp.FinishWithError("Failed to store stack definition.")
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}
	sp.FinishWithSuccess("Change hashes computed and stored successfully!")

	fmt.Println(utils.Success("All containers started successfully!"))
}
