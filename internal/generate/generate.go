package generate

import (
	"fmt"

	"github.com/danecwalker/podstack/internal/definition"
	"github.com/danecwalker/podstack/internal/spinners"
	"github.com/danecwalker/podstack/internal/utils"
)

type Generator interface {
	GenerateContainer(stack *definition.Stack, name string) ([]byte, error)
	GenerateNetwork(stack *definition.Stack, name string) ([]byte, error)
	GenerateVolume(stack *definition.Stack, name string) ([]byte, error)
}

func Generate(stack *definition.Stack, outputPath string, generator Generator) error {
	for name, network := range stack.Networks {
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
		sp := spinners.DefaultSpinner()
		container.ContainerName = name
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
