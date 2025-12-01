package definition

type Container struct {
	ContainerName string      `yaml:"-"`
	Entrypoint    StringArray `yaml:"entrypoint"`
	Environment   MapArray    `yaml:"environment"`
	// Healthcheck   *Healthcheck  `yaml:"healthcheck"`
	Image    *Image        `yaml:"image"`
	Init     bool          `yaml:"init"`
	Labels   MapArray      `yaml:"labels"`
	Networks []string      `yaml:"networks"`
	Ports    []PortMap     `yaml:"ports"`
	Restart  RestartPolicy `yaml:"restart"`
	Volumes  []VolumeMap   `yaml:"volumes"`
	Depends  []string      `yaml:"depends"`
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
