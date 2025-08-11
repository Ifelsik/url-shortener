package hasher

import (
	"hash"
	"hash/fnv"
)

type Hasher interface {
	String(hash string) string
}

type hasher32 struct {
	h hash.Hash
}

func NewHasher32() *hasher32 {
	return &hasher32{h: fnv.New32()}
}

func (h *hasher32) String(hash string) string {
	h.h.Reset()
	h.h.Write([]byte(hash))
	
	return string(h.h.Sum(nil))
}
