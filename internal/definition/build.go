package definition

import (
	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

type Build struct {
	Context       string   `yaml:"context"`
	ContainerFile string   `yaml:"containerfile"`
	Tags          []string `yaml:"tags"`
	Args          MapArray `yaml:"args"`
}

func (b *Build) UnmarshalYAML(value *yaml.Node) error {
	// If the build field is a string, treat it as the context
	if value.Kind == yaml.ScalarNode {
		var context string
		if err := value.Decode(&context); err != nil {
			return err
		}
		*b = Build{
			Context:       context,
			ContainerFile: "",
		}
		return nil
	}

	// Otherwise, treat it as a full Build struct
	type buildAlias Build
	var ba buildAlias
	if err := value.Decode(&ba); err != nil {
		return err
	}
	*b = Build(ba)
	return nil
}

func (b *Build) MarshalHash(h *hasher.Hash) error {
	if b == nil {
		return nil
	}
	h.Hasher.Write([]byte(b.Context))
	h.Hasher.Write([]byte(b.ContainerFile))
	for _, tag := range b.Tags {
		h.Hasher.Write([]byte(tag))
	}
	b.Args.MarshalHash(h)
	return nil
}
