package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerPort(t *testing.T) {
	tests := []struct {
		port  string
		valid *PortMap
	}{
		{"3000", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 3000,
				End:   3000,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 3000,
				End:   3000,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"3000-3005", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 3000,
				End:   3005,
				Range: true,
			},
			ContainerPort: &PortRange{
				Start: 3000,
				End:   3005,
				Range: true,
			},
			Protocol: "tcp",
		}},
		{"8000:8000", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 8000,
				End:   8000,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 8000,
				End:   8000,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"8080-8085:80-85/udp", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 8080,
				End:   8085,
				Range: true,
			},
			ContainerPort: &PortRange{
				Start: 80,
				End:   85,
				Range: true,
			},
			Protocol: "udp",
		}},
		{"9090-9091:8080-8081", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 9090,
				End:   9091,
				Range: true,
			},
			ContainerPort: &PortRange{
				Start: 8080,
				End:   8081,
				Range: true,
			},
			Protocol: "tcp",
		}},
		{"49100:22", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 49100,
				End:   49100,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 22,
				End:   22,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"8000-9000:80", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 8000,
				End:   9000,
				Range: true,
			},
			ContainerPort: &PortRange{
				Start: 80,
				End:   80,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"127.0.0.1:8001:8001", &PortMap{
			IP: "127.0.0.1",
			HostPort: &PortRange{
				Start: 8001,
				End:   8001,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 8001,
				End:   8001,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"127.0.0.1:5000-5010:5000-5010", &PortMap{
			IP: "127.0.0.1",
			HostPort: &PortRange{
				Start: 5000,
				End:   5010,
				Range: true,
			},
			ContainerPort: &PortRange{
				Start: 5000,
				End:   5010,
				Range: true,
			},
			Protocol: "tcp",
		}},
		{"[::1]:6000:6000", &PortMap{
			IP: "[::1]",
			HostPort: &PortRange{
				Start: 6000,
				End:   6000,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 6000,
				End:   6000,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"[::1]:6001:6001", &PortMap{
			IP: "[::1]",
			HostPort: &PortRange{
				Start: 6001,
				End:   6001,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 6001,
				End:   6001,
				Range: false,
			},
			Protocol: "tcp",
		}},
		{"6060:6060/udp", &PortMap{
			IP: "0.0.0.0",
			HostPort: &PortRange{
				Start: 6060,
				End:   6060,
				Range: false,
			},
			ContainerPort: &PortRange{
				Start: 6060,
				End:   6060,
				Range: false,
			},
			Protocol: "udp",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.port, func(t *testing.T) {
			p, err := ParsePort(tt.port)
			assert.NoError(t, err)
			assert.Equal(t, tt.valid, p)
		})
	}
}
