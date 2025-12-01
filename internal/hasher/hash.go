package hasher

import (
	"hash"

	"github.com/akamensky/base58"
	"lukechampine.com/blake3"
)

type Hash struct {
	Hasher hash.Hash
}

func NewHash() *Hash {
	return &Hash{
		Hasher: blake3.New(32, nil),
	}
}

type Hashable interface {
	MarshalHash(h *Hash) error
}

func MarshalHashable(obj Hashable) ([]byte, error) {
	h := NewHash()
	err := obj.MarshalHash(h)
	if err != nil {
		return nil, err
	}
	return h.Hasher.Sum(nil), nil
}

func MarshalHashableB58(obj Hashable) (string, error) {
	h := NewHash()
	err := obj.MarshalHash(h)
	if err != nil {
		return "", err
	}
	return EncodeB58(h.Hasher.Sum(nil)), nil
}

func EncodeB58(data []byte) string {
	base58Encoded := base58.Encode(data)
	return base58Encoded
}
