package msgpack

// http://godoc.org/github.com/youtube/vitess/go/rpcplus

import (
	"bufio"
	"io"

	"github.com/ugorji/go/codec"
	"github.com/youtube/vitess/go/rpcplus"
)

type PlusClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *codec.Decoder
	enc    *codec.Encoder
	encBuf *bufio.Writer
}

func NewPlusClientCodec(rwc io.ReadWriteCloser, mh *codec.MsgpackHandle) *PlusClientCodec {
	encBuf := bufio.NewWriter(rwc)
	return &PlusClientCodec{
		rwc:    rwc,
		dec:    codec.NewDecoder(rwc, mh),
		enc:    codec.NewEncoder(encBuf, mh),
		encBuf: encBuf,
	}
}

func (c *PlusClientCodec) WriteRequest(r *rpcplus.Request, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *PlusClientCodec) ReadResponseHeader(r *rpcplus.Response) error {
	return c.dec.Decode(r)
}

func (c *PlusClientCodec) ReadResponseBody(body interface{}) error {
	if body == nil {
		return c.dec.Decode(&body)
	}
	return c.dec.Decode(body)
}

func (c *PlusClientCodec) Close() error {
	return c.rwc.Close()
}

type PlusServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *codec.Decoder
	enc    *codec.Encoder
	encBuf *bufio.Writer
}

func NewPlusServerCodec(rwc io.ReadWriteCloser, mh *codec.MsgpackHandle) *PlusServerCodec {
	encBuf := bufio.NewWriter(rwc)
	return &PlusServerCodec{
		rwc:    rwc,
		dec:    codec.NewDecoder(rwc, mh),
		enc:    codec.NewEncoder(encBuf, mh),
		encBuf: encBuf,
	}
}

func (c *PlusServerCodec) ReadRequestHeader(r *rpcplus.Request) error {
	return c.dec.Decode(r)
}

func (c *PlusServerCodec) ReadRequestBody(body interface{}) error {
	if body == nil {
		return c.dec.Decode(&body)
	}
	return c.dec.Decode(body)
}

func (c *PlusServerCodec) WriteResponse(r *rpcplus.Response, body interface{}, last bool) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *PlusServerCodec) Close() error {
	return c.rwc.Close()
}
