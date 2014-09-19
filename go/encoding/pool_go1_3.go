// +build go1.3

package encoding

import (
	"bytes"
	"sync"
)

type EncoderPool struct {
	marshalers sync.Pool
}

func NewEncoderPool(factory EncoderFactory) *EncoderPool {
	t, _ := factory.(Trimer)
	return &EncoderPool{
		marshalers: sync.Pool{
			New: func() interface{} {
				NewEnc.Add(1)
				b := bytes.NewBuffer(nil)
				return marshaler{
					buffer:  b,
					encoder: factory.NewEncoder(b),
					trimer:  t,
				}
			},
		},
	}
}

func (p *EncoderPool) Marshal(v interface{}) ([]byte, error) {
	m := p.marshalers.Get().(marshaler)
	b, err := m.Marshal(v)
	p.marshalers.Put(m)
	return b, err
}

type DecoderPool struct {
	unmarshalers sync.Pool
}

func NewDecoderPool(factory DecoderFactory) *DecoderPool {
	return &DecoderPool{
		unmarshalers: sync.Pool{
			New: func() interface{} {
				b := bytes.NewBuffer(nil)
				return unmarshaler{
					buffer:  b,
					decoder: factory.NewDecoder(b),
				}
			},
		},
	}
}

func (p *DecoderPool) Unmarshal(b []byte, v interface{}) error {
	u := p.unmarshalers.Get().(unmarshaler)
	err := u.Unmarshal(b, v)
	p.unmarshalers.Put(u)
	return err
}
