package json

import (
	"encoding/json"
	"io"

	"git.gree-dev.net/twainy/tiroler-go/encoding"
)

type encoderFactory struct{}

func (f encoderFactory) NewEncoder(w io.Writer) encoding.Encoder {
	return json.NewEncoder(w)
}

func (f encoderFactory) Trim(b []byte) []byte {
	if len(b) > 0 {
		return b[:len(b)-1]
	}
	return b
}

var DefaultEncoderPool = NewEncoderPool()

func NewEncoderPool() *encoding.EncoderPool {
	return encoding.NewEncoderPool(encoderFactory{})
}

func Marshal(v interface{}) ([]byte, error) {
	return DefaultEncoderPool.Marshal(v)
}
