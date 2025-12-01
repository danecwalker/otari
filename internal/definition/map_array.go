package definition

import (
	"fmt"
	"sort"
	"strings"

	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

// MapArray is a custom type to handle YAML unmarshalling of map arrays.
//
// Values can be provided as a single map or as an array of KEY=VALUE strings.
// The output is always a single map with merged key-value pairs.
type MapArray map[string]string

func (ma *MapArray) UnmarshalYAML(node *yaml.Node) error {
	result := make(map[string]string)

	switch node.Kind {
	case yaml.MappingNode:
		if err := node.Decode(&result); err != nil {
			return fmt.Errorf("failed to decode mapping node: %w", err)
		}
	case yaml.SequenceNode:
		var items []string
		if err := node.Decode(&items); err != nil {
			return fmt.Errorf("failed to decode sequence node: %w", err)
		}
		for _, item := range items {
			parts := strings.SplitN(item, "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			}
		}
	default:
		return fmt.Errorf("unsupported YAML node kind for MapArray: %v", node.Kind)
	}

	*ma = result

	return nil
}

func (ma MapArray) MarshalHash(h *hasher.Hash) {
	// To ensure consistent hashing, sort the keys before processing.
	keys := make([]string, 0, len(ma))
	for k := range ma {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		h.Hasher.Write([]byte(k))
		h.Hasher.Write([]byte(ma[k]))
	}
}
