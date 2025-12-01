package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/danecwalker/podstack/internal/definition"
	"github.com/danecwalker/podstack/internal/generate"
	"github.com/danecwalker/podstack/internal/podman"
	"github.com/danecwalker/podstack/internal/quadlets"
	"github.com/danecwalker/podstack/internal/rules"
	"github.com/danecwalker/podstack/internal/spinners"
	"github.com/danecwalker/podstack/internal/systemd"
	"github.com/danecwalker/podstack/internal/utils"
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

	errors := rules.Validate(stack)
	if len(errors) > 0 {
		fmt.Println(utils.Error("Failed to validate stack:"))
		for _, err := range errors {
			color.New(color.FgRed).Printf("    • %s\n", err.Message)
		}
		os.Exit(1)
	}

	fmt.Println(utils.Success("Stack validated successfully!"))

	// check if containers
	if len(stack.Containers) == 0 {
		fmt.Println(utils.Info("No containers defined in the stack."))
		return
	}

	imagesSet := make(map[string]struct{})
	for _, container := range stack.Containers {
		imagesSet[container.Image.String()] = struct{}{}
	}

	for image := range imagesSet {
		sp := spinners.DefaultSpinner()

		// Check if image exists
		if !podman.ImageExists(ctx, image) {
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
			sp.FinishWithInfo(fmt.Sprintf("Image '%s' already exists.", image))
		}
	}

	// Generate systemd quadlets
	fmt.Println(utils.Info("Generating systemd quadlets..."))
	outputDir := utils.OutputLocation()
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println(utils.Error("Failed to create output directory."))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	if err := generate.Generate(stack, outputDir, quadlets.Generator()); err != nil {
		fmt.Println(utils.Error("Failed to generate systemd quadlets."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}

	time.Sleep(500 * time.Millisecond) // Small delay to ensure all spinners finish

	fmt.Println(utils.Success("Systemd quadlets generated successfully!"))
	fmt.Println()

	if err := systemd.ReloadDaemon(); err != nil {
		fmt.Println(utils.Error("Failed to reload systemd daemon."))
		color.New(color.FgWhite).Println("    " + err.Error())

		os.Exit(1)
	}

	// Start all containers
	for _, container := range stack.Containers {
		sp := spinners.DefaultSpinner()
		containerUnitName := container.ContainerName
		sp.SetMessage(fmt.Sprintf("Starting container '%s'...", containerUnitName))
		if err := systemd.StartUnit(containerUnitName); err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to start container '%s'", containerUnitName))
			color.New(color.FgWhite).Println("    " + err.Error())

			os.Exit(1)
		}
		sp.FinishWithSuccess(fmt.Sprintf("Container '%s' started.", containerUnitName))
	}

	fmt.Println(utils.Success("All containers started successfully!"))

	fmt.Printf("%+v\n", stack)
}
