package msgpack

import (
	"io"

	"github.com/twainy/tiroler-go/go/encoding"
	"github.com/ugorji/go/codec"
)

type decoderFactory struct {
	h *codec.MsgpackHandle
}

func (f decoderFactory) NewDecoder(r io.Reader) encoding.Decoder {
	return codec.NewDecoder(r, f.h)
}

var DefaultDecoderPool = NewDecoderPool(&codec.MsgpackHandle{})

func NewDecoderPool(h *codec.MsgpackHandle) *encoding.DecoderPool {
	return encoding.NewDecoderPool(decoderFactory{h: h})
}

func Unmarshal(b []byte, v interface{}) error {
	return DefaultDecoderPool.Unmarshal(b, v)
}
