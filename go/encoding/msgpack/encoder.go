package msgpack

import (
	"io"

	"github.com/twainy/tiroler-go/go/encoding"
	"github.com/ugorji/go/codec"
)

type encoderFactory struct {
	h *codec.MsgpackHandle
}

func (f encoderFactory) NewEncoder(w io.Writer) encoding.Encoder {
	return codec.NewEncoder(w, f.h)
}

var DefaultEncoderPool = NewEncoderPool(&codec.MsgpackHandle{})

func NewEncoderPool(h *codec.MsgpackHandle) *encoding.EncoderPool {
	return encoding.NewEncoderPool(encoderFactory{h: h})
}

func Marshal(v interface{}) ([]byte, error) {
	return DefaultEncoderPool.Marshal(v)
}
