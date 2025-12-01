package changes

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/danecwalker/otari/internal/definition"
	"github.com/danecwalker/otari/internal/hasher"
	"github.com/danecwalker/otari/internal/utils"
)

type StackData struct {
	Version     int               `toml:"version"`
	GeneratedAt time.Time         `toml:"generated_at"`
	Containers  map[string]string `toml:"containers,omitempty"`
	Volumes     map[string]string `toml:"volumes,omitempty"`
	Networks    map[string]string `toml:"networks,omitempty"`
}

func DetectChanges(ctx context.Context, newStack *definition.Stack) (new *definition.Stack, deleted *definition.Stack, total int, err error) {
	lockPath := newStack.StackName + ".lock"
	// check if lock file exists
	if !utils.PathExists(lockPath) {
		// no lock file, everything is new
		return newStack, nil, -1, nil
	}

	// read existing stack file
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil, nil, -1, err
	}

	var stackData StackData
	if err := toml.Unmarshal(data, &stackData); err != nil {
		return nil, nil, -1, err
	}

	if stackData.Version != 1 {
		// unsupported version, treat everything as new
		return newStack, nil, -1, fmt.Errorf("unsupported stack data version: %d", stackData.Version)
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
	stackData := StackData{
		Containers: make(map[string]string),
		Volumes:    make(map[string]string),
		Networks:   make(map[string]string),
	}
	for name, container := range stack.Containers {
		if container.Build != nil {
			container.Image = nil // do not hash build info
		}
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

	stackData.Version = 1
	stackData.GeneratedAt = time.Now().UTC().Truncate(time.Second)

	lockPath := stack.StackName + ".lock"
	f, err := os.Create(lockPath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := toml.NewEncoder(f)
	enc.Indent = ""
	err = enc.Encode(stackData)
	if err != nil {
		return err
	}

	return nil
}
