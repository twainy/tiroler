// +build !go1.3

package encoding

import (
	"bytes"
	"runtime"
)

type EncoderPool struct {
	marshalers chan marshaler
	factory    EncoderFactory
}

func NewEncoderPool(factory EncoderFactory) *EncoderPool {
	return &EncoderPool{
		marshalers: make(chan marshaler, runtime.NumCPU()*2),
		factory:    factory,
	}
}

func (p *EncoderPool) Marshal(v interface{}) ([]byte, error) {
	// get marshaler
	var m marshaler
	select {
	case m = <-p.marshalers:
	default:
		NewEnc.Add(1)
		m = marshaler{buffer: bytes.NewBuffer(nil)}
		m.encoder = p.factory.NewEncoder(m.buffer)
		m.trimer, _ = p.factory.(Trimer)
	}
	// marshal
	b, err := m.Marshal(v)
	// put marshaler
	select {
	case p.marshalers <- m:
	default:
	}
	return b, err
}

type DecoderPool struct {
	unmarshalers chan unmarshaler
	factory      DecoderFactory
}

func NewDecoderPool(factory DecoderFactory) *DecoderPool {
	return &DecoderPool{
		unmarshalers: make(chan unmarshaler, runtime.NumCPU()*2),
		factory:      factory,
	}
}

func (p *DecoderPool) Unmarshal(b []byte, v interface{}) error {
	// get unmarshaler
	var u unmarshaler
	select {
	case u = <-p.unmarshalers:
	default:
		b := bytes.NewBuffer(nil)
		u = unmarshaler{
			buffer:  b,
			decoder: p.factory.NewDecoder(b),
		}
	}
	// unmarshal
	err := u.Unmarshal(b, v)
	// put unmarshaler
	select {
	case p.unmarshalers <- u:
	default:
	}
	return err
}
