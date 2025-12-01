package quadlets

import (
	"github.com/danecwalker/podstack/internal/generate"
)

type QuadletGenerator struct {
	// Fields and methods for quadlet generation would go here
}

func Generator() generate.Generator {
	return &QuadletGenerator{}
}
