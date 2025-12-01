package definition

import (
	"fmt"

	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

type RestartPolicy struct {
	Condition   string `yaml:"condition"`
	MaxAttempts int    `yaml:"max_attempts"`
}

func (r *RestartPolicy) IsAlways() bool {
	return r.Condition == "always"
}

func (r *RestartPolicy) IsOnFailure() bool {
	return r.Condition == "on-failure"
}

func (r *RestartPolicy) IsNo() bool {
	return r.Condition == "no"
}

func (r *RestartPolicy) IsUnlessStopped() bool {
	return r.Condition == "unless-stopped"
}

func (r *RestartPolicy) UnmarshalYAML(value *yaml.Node) error {
	var policyStr string
	if err := value.Decode(&policyStr); err != nil {
		return err
	}

	switch policyStr {
	case "always":
		r.Condition = "always"
	case "no":
		r.Condition = "no"
	case "unless-stopped":
		r.Condition = "unless-stopped"
	default:
		// Check for on-failure with optional max attempts
		var maxAttempts int
		n, err := fmt.Sscanf(policyStr, "on-failure:%d", &maxAttempts)
		if err != nil || n == 0 {
			// Default to on-failure with no max attempts
			r.Condition = "on-failure"
			r.MaxAttempts = 0
		} else {
			r.Condition = "on-failure"
			r.MaxAttempts = maxAttempts
		}
	}

	return nil
}

func (r RestartPolicy) MarshalHash(h *hasher.Hash) {
	h.Hasher.Write([]byte(r.Condition))
	if r.Condition == "on-failure" {
		h.Hasher.Write([]byte(fmt.Sprintf(":%d", r.MaxAttempts)))
	}
}
