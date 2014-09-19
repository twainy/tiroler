package encoding

import (
	"bytes"
	"io"

	"github.com/youtube/vitess/go/sync2"
)

var NewEnc sync2.AtomicInt64

type Encoder interface {
	Encode(v interface{}) error
}

type EncoderFactory interface {
	NewEncoder(w io.Writer) Encoder
}

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
}

type Trimer interface {
	Trim(b []byte) []byte
}

type marshaler struct {
	buffer  *bytes.Buffer
	encoder Encoder
	trimer  Trimer
}

func (m marshaler) Marshal(v interface{}) ([]byte, error) {
	defer m.buffer.Reset()
	if err := m.encoder.Encode(v); err != nil {
		return nil, err
	}
	t := m.buffer.Bytes()
	if m.trimer != nil {
		t = m.trimer.Trim(t)
	}
	b := make([]byte, len(t))
	copy(b, t)
	return b, nil
}
