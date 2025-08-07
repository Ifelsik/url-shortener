package base62

import (
	"github.com/deatil/go-encoding/base62"
)

type Base62Provider interface {
	EncodeToString(data []byte) string
}

type base62Encoder struct {
	encoder *base62.Encoding
}

func NewBase62Encoder() *base62Encoder {
	return &base62Encoder{encoder: base62.StdEncoding}
}

func (e *base62Encoder) EncodeToString(data []byte) string {
	return e.encoder.EncodeToString(data)
}
