package hasher

import (
	"hash"
	"hash/fnv"
)

type Hasher interface {
	String(hash string) string
}

type hasher128 struct {
	h hash.Hash
}

func New128() *hasher128 {
	return &hasher128{h: fnv.New128()}
}

func (h *hasher128) String(hash string) string {
	h.h.Reset()
	h.h.Write([]byte(hash))
	return string(h.h.Sum(nil))
}
