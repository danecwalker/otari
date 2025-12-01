package definition

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// StringArray is a custom type to handle YAML unmarshalling of string arrays.
//
// Values can be provided as a single string or as an array of strings.
// The output is always a single space-separated string.
type StringArray string

func (sa *StringArray) UnmarshalYAML(node *yaml.Node) error {
	var result []string

	switch node.Kind {
	case yaml.ScalarNode:
		result = []string{node.Value}
	case yaml.SequenceNode:
		if err := node.Decode(&result); err != nil {
			return err
		}
	default:
		return nil
	}

	*sa = StringArray(strings.Join(result, " "))

	return nil
}
