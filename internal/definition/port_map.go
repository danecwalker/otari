package definition

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

var portRe = regexp.MustCompile(
	`^(?:((?:\d{1,3}(?:\.\d{1,3}){3})|\[.+?\]))?:?` + // (1) IP (optional)
		`(` +
		`\d+(?:-\d+)?` + // (2) Host port OR single port
		`)` +
		`(?:` +
		`:` +
		`(\d+)(?:-(\d+))?` + // (3),(4) Container port / optional range
		`)?` +
		`(?:/(tcp|udp))?$`, // (5) Protocol
)

type PortMap struct {
	IP            string
	HostPort      *PortRange
	ContainerPort *PortRange
	Protocol      string
}

type PortRange struct {
	Start int
	End   int
	Range bool
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func parseRange(s string) (start, end int, isRange bool) {
	if strings.Contains(s, "-") {
		parts := strings.SplitN(s, "-", 2)
		return atoi(parts[0]), atoi(parts[1]), true
	}
	v := atoi(s)
	return v, v, false
}

func ParsePort(portStr string) (*PortMap, error) {
	m := portRe.FindStringSubmatch(portStr)
	if m == nil {
		return nil, fmt.Errorf("invalid port format: %s", portStr)
	}

	ip := m[1]
	if ip == "" {
		ip = "0.0.0.0"
	}

	// m[2]: host port or single port
	hStart, hEnd, hRange := parseRange(m[2])

	var cStart, cEnd int
	var cRange bool

	if m[3] == "" {
		// No explicit container port → use host ports for container too
		cStart, cEnd, cRange = hStart, hEnd, hRange
	} else {
		// m[3]: container start, m[4]: container end (optional)
		if m[4] != "" {
			cStart, cEnd, cRange = parseRange(m[3] + "-" + m[4])
		} else {
			cStart, cEnd, cRange = parseRange(m[3])
		}
	}

	protocol := m[5]
	if protocol == "" {
		protocol = "tcp"
	}

	return &PortMap{
		IP: ip,
		HostPort: &PortRange{
			Start: hStart,
			End:   hEnd,
			Range: hRange,
		},
		ContainerPort: &PortRange{
			Start: cStart,
			End:   cEnd,
			Range: cRange,
		},
		Protocol: protocol,
	}, nil
}

func (p *PortMap) UnmarshalYAML(value *yaml.Node) error {
	var portStr string
	if err := value.Decode(&portStr); err != nil {
		return err
	}

	parsedPort, err := ParsePort(portStr)
	if err != nil {
		return err
	}

	*p = *parsedPort
	return nil
}

func (p *PortMap) String() string {
	var sb strings.Builder

	sb.WriteString(p.IP)
	sb.WriteString(":")

	if p.HostPort.Range {
		sb.WriteString(fmt.Sprintf("%d-%d", p.HostPort.Start, p.HostPort.End))
	} else {
		sb.WriteString(fmt.Sprintf("%d", p.HostPort.Start))
	}

	sb.WriteString(":")

	if p.ContainerPort.Range {
		sb.WriteString(fmt.Sprintf("%d-%d", p.ContainerPort.Start, p.ContainerPort.End))
	} else {
		sb.WriteString(fmt.Sprintf("%d", p.ContainerPort.Start))
	}

	sb.WriteString("/")
	sb.WriteString(p.Protocol)

	return sb.String()
}

func (p PortMap) MarshalHash(h *hasher.Hash) {
	h.Hasher.Write([]byte(p.String()))
}
