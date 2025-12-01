package generate

import (
	"fmt"

	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/spinners"
	"github.com/danecwalker/otari/internal/utils"
)

type Generator interface {
	GenerateContainer(stack *definition.Stack, name string) ([]byte, error)
	GenerateNetwork(stack *definition.Stack, name string) ([]byte, error)
	GenerateVolume(stack *definition.Stack, name string) ([]byte, error)
}

func Generate(stack, new *definition.Stack, outputPath string, generator Generator) error {
	for name, network := range stack.Networks {
		if _, exists := new.Networks[name]; !exists {
			continue
		}
		if network.Driver == definition.NetworkDriverHost {
			continue
		}
		sp := spinners.DefaultSpinner()
		network.NetworkName = name
		sp.SetMessage(fmt.Sprintf("Generating configuration for network '%s'", network.NetworkName))
		out, err := generator.GenerateNetwork(stack, network.NetworkName)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for network '%s': %v", network.NetworkName, err))
			return err
		}

		err = utils.WriteToFile(outputPath, network.NetworkName+".network", out)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for network '%s': %v", network.NetworkName, err))
			return err
		}

		sp.FinishWithSuccess(fmt.Sprintf("Generated configuration for network '%s'", network.NetworkName))
	}

	for name, volume := range stack.Volumes {
		if _, exists := new.Volumes[name]; !exists {
			continue
		}
		sp := spinners.DefaultSpinner()
		volume.VolumeName = name
		sp.SetMessage(fmt.Sprintf("Generating configuration for volume '%s'", volume.VolumeName))
		out, err := generator.GenerateVolume(stack, volume.VolumeName)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for volume '%s': %v", volume.VolumeName, err))
			return err
		}

		err = utils.WriteToFile(outputPath, volume.VolumeName+".volume", out)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for volume '%s': %v", volume.VolumeName, err))
			return err
		}
		sp.FinishWithSuccess(fmt.Sprintf("Generated configuration for volume '%s'", volume.VolumeName))
	}

	for name, container := range stack.Containers {
		if _, exists := new.Containers[name]; !exists {
			continue
		}
		sp := spinners.DefaultSpinner()
		container.ContainerName = name
		if container.Build != nil {
			container.Image.Registry = "localhost"
			container.Image.Image = fmt.Sprintf("%s_%s", stack.StackName, container.ContainerName)
		}
		sp.SetMessage(fmt.Sprintf("Generating configuration for container '%s'", container.ContainerName))
		out, err := generator.GenerateContainer(stack, container.ContainerName)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for container '%s': %v", container.ContainerName, err))
			return err
		}

		err = utils.WriteToFile(outputPath, container.ContainerName+".container", out)
		if err != nil {
			sp.FinishWithError(fmt.Sprintf("Failed to generate configuration for container '%s': %v", container.ContainerName, err))
			return err
		}
		sp.FinishWithSuccess(fmt.Sprintf("Generated configuration for container '%s'", container.ContainerName))
	}
	return nil
}
