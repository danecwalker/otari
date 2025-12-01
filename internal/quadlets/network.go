package quadlets

import (
	"bytes"
	"fmt"

	"github.com/danecwalker/podstack/internal/definition"
	"github.com/danecwalker/podstack/internal/utils"
)

func (q *QuadletGenerator) GenerateNetwork(stack *definition.Stack, networkName string) ([]byte, error) {
	network, exists := stack.Networks[networkName]
	if !exists {
		return nil, fmt.Errorf("network '%s' not found in stack", networkName)
	}
	var buf bytes.Buffer
	err := utils.WriteSection(&buf, "Unit", [][2]string{
		{"Description", fmt.Sprintf("%s network", network.NetworkName)},
	})
	if err != nil {
		return nil, err
	}

	if err := utils.WriteEmptyLine(&buf); err != nil {
		return nil, err
	}

	networkProperties := [][2]string{
		{"NetworkName", network.NetworkName},
	}

	if network.Driver != definition.NetworkDriverHost && network.Driver != "" {
		networkProperties = append(networkProperties, [2]string{
			"Driver", string(network.Driver),
		})
	}

	err = utils.WriteSection(&buf, "Network", networkProperties)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
