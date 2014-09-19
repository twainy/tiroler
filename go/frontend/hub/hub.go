package hub

import (
	"sync"

	"git.gree-dev.net/giistudio/blackbird-go/go/encoding/json"
	"github.com/twainy/tiroler-go/socket"
	"github.com/golang/glog"
)

const (
	trace = true
	test  = true
)

var ejectMsg = []byte(`{"type":"Connection.Eject"}`)

type typer struct {
	Type string `json:"type"`
}

type handler func(*conn, []byte) error

type Hub struct {
	mu       sync.Mutex
	conns    map[socket.ID]*conn
	auths    map[uint64]*conn
	handlers map[string]handler
}

func NewHub() *Hub {
	h := &Hub{
		conns: make(map[socket.ID]*conn),
		auths: make(map[uint64]*conn),
		handlers: map[string]handler{
			rpcRequestType: rpcRequest,
			timePingType:   timePing,
		},
	}
	if test {
		h.handlers[testEchoType] = testEcho
		h.handlers[testDelayType] = testDelay
	}
	return h
}

func (h *Hub) Connect(c socket.Conn) {
	conn := newConn(c)
	if trace {
		glog.Infoln(conn.LogTag(), "connect")
	}
	h.mu.Lock()
	h.conns[c.ID()] = conn
	h.mu.Unlock()
	go h.listen(conn)
}

func (h *Hub) Disconnect(c socket.Conn) {
	h.mu.Lock()
	if conn, ok := h.conns[c.ID()]; ok {
		if trace {
			glog.Infoln(conn.LogTag(), "disconnect")
		}
		delete(h.conns, conn.ID())
		delete(h.auths, conn.AuthID())
		c.Close(nil)
	}
	h.mu.Unlock()
}

func (h *Hub) listen(c *conn) {
	defer h.Disconnect(c)
	for data := range c.Receive() {
		var t typer
		if err := json.Unmarshal(data, &t); err != nil {
			glog.Errorln(c.LogTag(), "unmarshal", err)
			return
		}
		handler, ok := h.handlers[t.Type]
		if !ok {
			glog.Errorln(c.LogTag(), "unhandled type", string(data))
			return
		}
		if err := handler(c, data); err != nil {
			return
		}
	}
}

func (h *Hub) auth(c *conn) {
	if trace {
		glog.Infoln(c.LogTag(), "auth")
	}
	id := c.AuthID()
	h.mu.Lock()
	if a, ok := h.auths[id]; ok {
		if trace {
			glog.Infoln(a.LogTag(), "eject")
		}
		a.Close(ejectMsg)
		delete(h.conns, a.ID())
	}
	h.auths[id] = c
	h.mu.Unlock()
}

var DefaultHub = NewHub()

func Connect(c socket.Conn) {
	DefaultHub.Connect(c)
}

func Disconnect(c socket.Conn) {
	DefaultHub.Disconnect(c)
}
