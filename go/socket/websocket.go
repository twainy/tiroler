package socket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const (
	trace = true

	sendBufferSize = 64              // Maximum number of writes to buffer
	maxMsgSize     = 1024 * 1024 * 4 // Maximum message size allowed from client.

	writeWait = 5 * time.Second // Time allowed to write a message to the client.

	pingPeriod = 45 * time.Second         // Send pings to peer with this period.
	pongWait   = pingPeriod + writeWait*2 // Time allowed to read the next pong message from the peer.
)

type WebSocket struct {
	conn *websocket.Conn

	id         ID
	remoteAddr string
	requestURI string

	recv     chan []byte
	send     chan []byte
	close    sync.Once
	closeMsg []byte
	closing  chan struct{}
	closed   chan struct{}
}

func NewWebSocket(w http.ResponseWriter, r *http.Request) *WebSocket {
	if r.Method != "GET" {
		glog.Errorln("method", r.RequestURI, r.RemoteAddr, r.Method)
		http.Error(w, "Method not allowed", 405)
		return nil
	}
	// if origin, err := url.Parse(r.Header.Get("Origin")); err != nil || origin.Host != r.Host {
	// 	glog.Errorln("origin not allowed", r.RequestURI, r.RemoteAddr, r.Header.Get("Origin"))
	// 	http.Error(w, "Origin not allowed", 403)
	// 	return nil
	// }
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		glog.Errorln("handshake", r.RequestURI, r.RemoteAddr)
		http.Error(w, "Not a websocket handshake", 400)
		return nil
	} else if err != nil {
		glog.Errorln("upgrade", r.RequestURI, r.RemoteAddr, err)
		return nil
	}
	ws := &WebSocket{
		conn:       conn,
		id:         ID(id.Add(1)),
		remoteAddr: r.RemoteAddr,
		requestURI: r.RequestURI,
		recv:       make(chan []byte),
		send:       make(chan []byte, sendBufferSize),
		closing:    make(chan struct{}),
		closed:     make(chan struct{}),
	}
	if trace {
		glog.Infoln(ws.LogTag(), "connect")
	}
	go ws.readPump()
	go ws.writePump()
	return ws
}

func (ws *WebSocket) readPump() {
	defer close(ws.recv)
	ws.conn.SetReadLimit(maxMsgSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error {
		if trace {
			glog.Infoln(ws.LogTag(), "pong")
		}
		ws.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, data, err := ws.conn.ReadMessage()
		if err != nil {
			ws.Close(nil)
			return
		}
		if trace {
			glog.Infoln(ws.LogTag(), "receive", string(data))
		}
		select {
		case ws.recv <- data:
			// continue pumping
		case <-ws.closing:
			return // read failed
		}
	}
}

func (ws *WebSocket) writePump() {
	pinger := time.NewTicker(pingPeriod)
	defer func() {
		pinger.Stop()
		ws.conn.Close()
		if trace {
			glog.Infoln(ws.LogTag(), "disconnect")
		}
		close(ws.closed)
	}()
	for {
		select {
		case msg := <-ws.send:
			if trace {
				glog.Infoln(ws.LogTag(), "send", string(msg))
			}
			if ws.write(websocket.BinaryMessage, msg) != nil {
				ws.Close(nil)
				return // write failed
			}
		case <-pinger.C:
			if trace {
				glog.Infoln(ws.LogTag(), "ping")
			}
			if ws.write(websocket.PingMessage, []byte{}) != nil {
				ws.Close(nil)
				return // write failed
			}
		case <-ws.closing:
			if ws.closeMsg != nil {
				if trace {
					glog.Infoln(ws.LogTag(), "data", string(ws.closeMsg))
				}
				ws.write(websocket.BinaryMessage, ws.closeMsg)
			}
			if trace {
				glog.Infoln(ws.LogTag(), "close")
			}
			ws.write(websocket.CloseMessage, []byte{})
			return // closing
		}
	}
}

func (ws *WebSocket) write(mt int, msg []byte) error {
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, msg)
}

func (ws *WebSocket) ID() ID {
	return ws.id
}

func (ws *WebSocket) RemoteAddr() string {
	return ws.remoteAddr
}

func (ws *WebSocket) Send(p []byte) error {
	select {
	case <-ws.closing:
		return ErrClosing
	default:
		select {
		case ws.send <- p:
			return nil
		default:
			return ErrBufferFull
		}
	}
}

func (ws *WebSocket) Receive() <-chan []byte {
	return ws.recv
}

func (ws *WebSocket) Close(p []byte) {
	ws.close.Do(func() {
		if trace {
			glog.Infoln(ws.LogTag(), "closing", string(p))
		}
		ws.closeMsg = p
		close(ws.closing)
	})
}

func (ws *WebSocket) Closing() <-chan struct{} {
	return ws.closing
}

func (ws *WebSocket) Closed() <-chan struct{} {
	return ws.closed
}

func (ws *WebSocket) LogTag() string {
	return fmt.Sprintf("%s %s %d", ws.requestURI, ws.remoteAddr, ws.id)
}
