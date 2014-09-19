package hub

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"git.gree-dev.net/giistudio/blackbird-go/go/socket"
	"github.com/youtube/vitess/go/sync2"
)

type conn struct {
	socket.Conn
	authID AtomicUint64

	rpcLock  sync.Mutex
	seq      sync2.AtomicUint32
	session  *json.RawMessage
	requests []*rpcClientRequest
}

func newConn(c socket.Conn) *conn {
	return &conn{
		Conn:     c,
		requests: make([]*rpcClientRequest, 0),
	}
}

func (c *conn) AuthID() uint64 {
	return c.authID.Get()
}

func (c *conn) LogTag() string {
	auth := c.AuthID()
	if auth == 0 {
		return fmt.Sprintf("%s ???", c.Conn.LogTag())
	}
	return fmt.Sprintf("%s %d", c.Conn.LogTag(), auth)
}

type AtomicUint64 uint64

func (i *AtomicUint64) Add(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Set(n uint64) {
	atomic.StoreUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Get() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}

func (i *AtomicUint64) CompareAndSwap(oldval, newval uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64((*uint64)(i), oldval, newval)
}
