package definition

import (
	"github.com/danecwalker/otari/internal/hasher"
	"github.com/danecwalker/otari/internal/systemd"
)

type Container struct {
	ContainerName string      `yaml:"-"`
	Entrypoint    StringArray `yaml:"entrypoint"`
	Environment   MapArray    `yaml:"environment"`
	// Healthcheck   *Healthcheck  `yaml:"healthcheck"`
	Image         *Image        `yaml:"image"`
	Build         *Build        `yaml:"build"`
	Init          bool          `yaml:"init"`
	Labels        MapArray      `yaml:"labels"`
	Networks      []string      `yaml:"networks"`
	Ports         []PortMap     `yaml:"ports"`
	RestartPolicy RestartPolicy `yaml:"restart"`
	Volumes       []VolumeMap   `yaml:"volumes"`
	Depends       []string      `yaml:"depends"`
}

// type Healthcheck struct {
// 	Test          StringArray   `yaml:"test"`
// 	Interval      time.Duration `yaml:"interval"`
// 	Timeout       time.Duration `yaml:"timeout"`
// 	Retries       int           `yaml:"retries"`
// 	StartPeriod   time.Duration `yaml:"start_period"`
// 	StartInterval time.Duration `yaml:"start_interval"`
// 	Disable       bool          `yaml:"disable"`
// }

func (c *Container) MarshalHash(h *hasher.Hash) error {
	if c == nil {
		return nil
	}
	h.Hasher.Write([]byte(c.ContainerName))
	c.Entrypoint.MarshalHash(h)
	c.Environment.MarshalHash(h)
	// if c.Healthcheck != nil {
	// 	c.Healthcheck.MarshalHash(h)
	// }
	if c.Image != nil {
		c.Image.MarshalHash(h)
	}
	if c.Build != nil {
		c.Build.MarshalHash(h)
	}
	if c.Init {
		h.Hasher.Write([]byte{1})
	} else {
		h.Hasher.Write([]byte{0})
	}
	c.Labels.MarshalHash(h)
	for _, net := range c.Networks {
		h.Hasher.Write([]byte(net))
	}
	for _, port := range c.Ports {
		port.MarshalHash(h)
	}
	c.RestartPolicy.MarshalHash(h)
	for _, vol := range c.Volumes {
		vol.MarshalHash(h)
	}
	for _, dep := range c.Depends {
		h.Hasher.Write([]byte(dep))
	}
	return nil
}

func (c *Container) Start() error {
	return systemd.StartUnit(c.ContainerName)
}

func (c *Container) Stop() error {
	return systemd.StopUnit(c.ContainerName)
}

func (c *Container) Restart() error {
	return systemd.RestartUnit(c.ContainerName)
}

func (c *Container) Remove() error {
	// Stop the container unit
	if err := systemd.StopUnit(c.ContainerName); err != nil {
		return err
	}

	// Remove the container quadlet
	if err := systemd.DeleteUnitFile(c.ContainerName + ".container"); err != nil {
		return err
	}

	return nil
}
