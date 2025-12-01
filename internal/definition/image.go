package definition

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/danecwalker/otari/internal/hasher"
	"gopkg.in/yaml.v3"
)

var imageRefRe = regexp.MustCompile(
	`^(?:` +
		`([a-zA-Z0-9.-]+(?::[0-9]+)?)/` + // 1: registry (optional)
		`)?` +
		`(` +
		`[a-z0-9]+(?:[._-][a-z0-9]+)*` + // first path component
		`(?:/[a-z0-9]+(?:[._-][a-z0-9]+)*)*` + // optional /more/components
		`)` + // 2: full path (project + image)
		`(?:` +
		`:([\w][\w.-]{0,127})` + // 3: tag (optional)
		`|` +
		`@(sha256:[0-9a-f]{64})` + // 4: digest (optional)
		`)?$`,
)

type Image struct {
	Registry  string
	Project   string
	Image     string
	Tag       string
	Digest    string
	FullyQual bool
}

func (img *Image) IsFullyQualified() bool {
	return img.FullyQual
}

func (img *Image) String() string {
	var sb strings.Builder
	if img.Registry != "" {
		sb.WriteString(img.Registry)
		sb.WriteString("/")
	}
	if img.Project != "" {
		sb.WriteString(img.Project)
		sb.WriteString("/")
	}
	sb.WriteString(img.Image)
	if img.Tag != "" {
		sb.WriteString(":")
		sb.WriteString(img.Tag)
	} else if img.Digest != "" {
		sb.WriteString("@")
		sb.WriteString(img.Digest)
	}
	return sb.String()
}

func ParseImage(imageRef string) (*Image, error) {
	m := imageRefRe.FindStringSubmatch(imageRef)
	if m == nil {
		return nil, fmt.Errorf("invalid image reference: %s", imageRef)
	}

	registry := m[1]
	fullPath := m[2] // always non-empty
	parts := strings.Split(fullPath, "/")
	var project, imageName string
	if len(parts) == 1 {
		imageName = parts[0]
	} else {
		project = strings.Join(parts[:len(parts)-1], "/")
		imageName = parts[len(parts)-1]
	}

	tag := m[3]
	digest := m[4]

	fullyQual := registry != "" || project != ""
	return &Image{
		Registry:  registry,
		Project:   project,
		Image:     imageName,
		Tag:       tag,
		Digest:    digest,
		FullyQual: fullyQual,
	}, nil
}

func (img *Image) UnmarshalYAML(node *yaml.Node) error {
	var imageRef string
	if err := node.Decode(&imageRef); err != nil {
		return err
	}

	parsedImage, err := ParseImage(imageRef)
	if err != nil {
		return err
	}

	*img = *parsedImage

	return nil
}

func (img Image) MarshalHash(h *hasher.Hash) {
	h.Hasher.Write([]byte(img.String()))
}
