package definition

import "github.com/danecwalker/otari/internal/hasher"

type Volume struct {
	VolumeName      string `yaml:"-"`
	PersistOnRemove bool   `yaml:"persist_on_remove"`
}

func (v *Volume) MarshalHash(h *hasher.Hash) error {
	if v == nil {
		return nil
	}
	h.Hasher.Write([]byte(v.VolumeName))
	return nil
}
