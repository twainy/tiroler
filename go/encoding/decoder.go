package encoding

import (
	"bytes"
	"io"
)

type Decoder interface {
	Decode(v interface{}) error
}

type DecoderFactory interface {
	NewDecoder(r io.Reader) Decoder
}

type Unmarshaler interface {
	Unmarshal(b []byte, v interface{}) error
}

type unmarshaler struct {
	buffer  *bytes.Buffer
	decoder Decoder
}

func (e unmarshaler) Unmarshal(b []byte, v interface{}) error {
	e.buffer.Reset()
	e.buffer.Write(b)
	return e.decoder.Decode(v)
}
