package definition

import (
	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

type NetworkDriver string

const (
	NetworkDriverBridge  NetworkDriver = "bridge"
	NetworkDriverHost    NetworkDriver = "host"
	NetworkDriverOverlay NetworkDriver = "ipvlan"
	NetworkDriverMacvlan NetworkDriver = "macvlan"
	NetworkDriverDefault NetworkDriver = NetworkDriverBridge
)

type Network struct {
	NetworkName     string        `yaml:"-"`
	Driver          NetworkDriver `yaml:"driver"`
	PersistOnRemove bool          `yaml:"persist_on_remove"`
}

func (n *NetworkDriver) UnmarshalYAML(value *yaml.Node) error {
	var driverStr string
	if err := value.Decode(&driverStr); err != nil {
		return err
	}

	switch driverStr {
	case "bridge":
		*n = NetworkDriverBridge
	case "host":
		*n = NetworkDriverHost
	case "ipvlan":
		*n = NetworkDriverOverlay
	case "macvlan":
		*n = NetworkDriverMacvlan
	default:
		*n = NetworkDriverDefault
	}

	return nil
}

func (n *Network) MarshalHash(h *hasher.Hash) error {
	if n == nil {
		return nil
	}
	h.Hasher.Write([]byte(n.NetworkName))
	h.Hasher.Write([]byte(n.Driver))
	return nil
}
