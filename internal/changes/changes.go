package changes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/hasher"
	"github.com/danecwalker/otari/internal/utils"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

type StackData struct {
	Containers map[string]string `yaml:"containers"`
	Volumes    map[string]string `yaml:"volumes"`
	Networks   map[string]string `yaml:"networks"`
}

func DetectChanges(ctx context.Context, newStack *definition.Stack) (new *definition.Stack, deleted *definition.Stack, total int, err error) {
	// check if stack exists in data dir
	dataDir := utils.DataDirectory()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Println(utils.Error("Failed to create data directory."))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	stackFilePath := filepath.Join(dataDir, newStack.StackName+".yaml")
	_, err = os.Stat(stackFilePath)
	if os.IsNotExist(err) {
		// stack does not exist, everything is new
		return newStack, nil, -1, nil
	} else if err != nil {
		return nil, nil, -1, err
	}

	// read existing stack file
	data, err := os.ReadFile(stackFilePath)
	if err != nil {
		return nil, nil, -1, err
	}

	var stackData StackData
	if err := yaml.Unmarshal(data, &stackData); err != nil {
		return nil, nil, -1, err
	}

	// stack data contains hashes of existing resources
	existingContainers := stackData.Containers
	existingVolumes := stackData.Volumes
	existingNetworks := stackData.Networks

	new = &definition.Stack{
		Containers: make(map[string]*definition.Container),
		Volumes:    make(map[string]*definition.Volume),
		Networks:   make(map[string]*definition.Network),
	}
	deleted = &definition.Stack{
		Containers: make(map[string]*definition.Container),
		Volumes:    make(map[string]*definition.Volume),
		Networks:   make(map[string]*definition.Network),
	}

	// detect new and modified containers
	for name, container := range newStack.Containers {
		hash, err := hasher.MarshalHashableB58(container)
		if err != nil {
			return nil, nil, -1, err
		}
		if existingHash, ok := existingContainers[name]; !ok || existingHash != hash {
			new.Containers[name] = container
		}
	}
	// detect deleted containers
	for name := range existingContainers {
		if _, ok := newStack.Containers[name]; !ok {
			deleted.Containers[name] = &definition.Container{ContainerName: name}
		}
	}
	// detect new and modified volumes
	for name, volume := range newStack.Volumes {
		hash, err := hasher.MarshalHashableB58(volume)
		if err != nil {
			return nil, nil, -1, err
		}
		if existingHash, ok := existingVolumes[name]; !ok || existingHash != hash {
			new.Volumes[name] = volume
		}
	}
	// detect deleted volumes
	for name := range existingVolumes {
		if _, ok := newStack.Volumes[name]; !ok {
			deleted.Volumes[name] = &definition.Volume{VolumeName: name}
		}
	}
	// detect new and modified networks
	for name, network := range newStack.Networks {
		hash, err := hasher.MarshalHashableB58(network)
		if err != nil {
			return nil, nil, -1, err
		}
		if existingHash, ok := existingNetworks[name]; !ok || existingHash != hash {
			new.Networks[name] = network
		}
	}
	// detect deleted networks
	for name := range existingNetworks {
		if _, ok := newStack.Networks[name]; !ok {
			deleted.Networks[name] = &definition.Network{NetworkName: name}
		}
	}

	totalChanges := len(new.Containers) + len(deleted.Containers) + len(new.Volumes) + len(deleted.Volumes) +
		len(new.Networks) + len(deleted.Networks)

	return new, deleted, totalChanges, nil
}

func SaveStackData(stack *definition.Stack) error {
	dataDir := utils.DataDirectory()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Println(utils.Error("Failed to create data directory."))
		color.New(color.FgWhite).Println("    " + err.Error())
		os.Exit(1)
	}

	stackFilePath := filepath.Join(dataDir, stack.StackName+".yaml")
	stackData := StackData{
		Containers: make(map[string]string),
		Volumes:    make(map[string]string),
		Networks:   make(map[string]string),
	}
	for name, container := range stack.Containers {
		hash, err := hasher.MarshalHashableB58(container)
		if err != nil {
			return err
		}
		stackData.Containers[name] = hash
	}
	for name, volume := range stack.Volumes {
		hash, err := hasher.MarshalHashableB58(volume)
		if err != nil {
			return err
		}
		stackData.Volumes[name] = hash
	}
	for name, network := range stack.Networks {
		hash, err := hasher.MarshalHashableB58(network)
		if err != nil {
			return err
		}
		stackData.Networks[name] = hash
	}
	data, err := yaml.Marshal(&stackData)
	if err != nil {
		return err
	}
	if err := os.WriteFile(stackFilePath, data, 0644); err != nil {
		return err
	}
	return nil
}
