package socket

import (
	"errors"

	"github.com/youtube/vitess/go/sync2"
)

type ID uint32

var (
	id            sync2.AtomicUint32
	ErrClosing    = errors.New("socket/conn: connection closing")
	ErrBufferFull = errors.New("socket/conn: buffer is full")
)

type Conn interface {
	// ID returns a unique identifier (scoped to the current server) for the connection.
	ID() ID
	// RemoteAddr returns the remote address of the client.
	RemoteAddr() string
	// Send sends a message to the client.
	Send(p []byte) error
	// Receive returns a channel of messages received from client. The channel is closed when the
	// connection is shutting down.
	Receive() <-chan []byte
	// Close signals the connection to shut down with an optional final message to send.
	Close(p []byte)
	// Closing returns a channel that signals once the connection has started closing.
	Closing() <-chan struct{}
	// Closed returns a channel that signals when connection has completed closing.
	Closed() <-chan struct{}

	LogTag() string
}
