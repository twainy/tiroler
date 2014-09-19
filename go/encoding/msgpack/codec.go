package msgpack

// http://godoc.org/github.com/youtube/vitess/go/rpcplus

import (
	"bufio"
	"io"
	"net/rpc"

	"github.com/ugorji/go/codec"
)

type ClientCodec struct {
	rwc    io.ReadWriteCloser
	dec    *codec.Decoder
	enc    *codec.Encoder
	encBuf *bufio.Writer
}

func NewClientCodec(rwc io.ReadWriteCloser, mh *codec.MsgpackHandle) *ClientCodec {
	encBuf := bufio.NewWriter(rwc)
	return &ClientCodec{
		rwc:    rwc,
		dec:    codec.NewDecoder(rwc, mh),
		enc:    codec.NewEncoder(encBuf, mh),
		encBuf: encBuf,
	}
}

func (c *ClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *ClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *ClientCodec) ReadResponseBody(body interface{}) error {
	if body == nil {
		return c.dec.Decode(&body)
	}
	return c.dec.Decode(body)
}

func (c *ClientCodec) Close() error {
	return c.rwc.Close()
}

type ServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *codec.Decoder
	enc    *codec.Encoder
	encBuf *bufio.Writer
}

func NewServerCodec(rwc io.ReadWriteCloser, mh *codec.MsgpackHandle) *ServerCodec {
	encBuf := bufio.NewWriter(rwc)
	return &ServerCodec{
		rwc:    rwc,
		dec:    codec.NewDecoder(rwc, mh),
		enc:    codec.NewEncoder(encBuf, mh),
		encBuf: encBuf,
	}
}

func (c *ServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.dec.Decode(r)
}

func (c *ServerCodec) ReadRequestBody(body interface{}) error {
	if body == nil {
		return c.dec.Decode(&body)
	}
	return c.dec.Decode(body)
}

func (c *ServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = c.enc.Encode(r); err != nil {
		return
	}
	if err = c.enc.Encode(body); err != nil {
		return
	}
	return c.encBuf.Flush()
}

func (c *ServerCodec) Close() error {
	return c.rwc.Close()
}
