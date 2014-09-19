package logging

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

const atomicOps = true // compiler will optimize away branches on constants and dead code

type handler struct {
	http.Handler
}

func Handler(h http.Handler) http.Handler {
	return handler{Handler: h}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	// WARNING: this assumes that we're sitting behind ELB and are not accessed directly
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For = clientIp, firstProxyIp, secondProxyIp
		// r.RemoteAddr = lastProxyIp:port
		if i := strings.LastIndex(xff, ", "); i != -1 {
			// if the request has gone through multiple proxies take the last ip
			// any intermediate ip can easily be forged by setting an x-forwarded-for
			// header by the client or an intermediate proxy
			r.RemoteAddr = xff[i+2:]
		} else {
			r.RemoteAddr = xff
		}
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		r.RemoteAddr = host
	}
	var logger responseWriter
	if _, ok := w.(http.Hijacker); ok {
		logger = &hijackLogger{responseLogger: &responseLogger{w: w}}
	} else {
		logger = &responseLogger{w: w}
	}
	h.Handler.ServeHTTP(logger, r)
	s := time.Since(t).Seconds() // connection duration seconds
	if logger.Hijacked() {
		// the only thing that we're hijacking connections with are websockets
		u, d := logger.Recv(), logger.Sent() // bytes upstream, bytes downstream
		ru, rd := float64(u)/s, float64(d)/s // bytes upstream/second, bytes downstream/second
		glog.Infof("ws %s %s %d %f %d %f %f\n", r.URL.RequestURI(), r.RemoteAddr, u, ru, d, rd, s)
	} else {
		glog.Infof("req %s %s %d %d %f\n", r.URL.RequestURI(), r.RemoteAddr, logger.Status(), logger.Sent(), s)
	}
}

type responseWriter interface {
	http.ResponseWriter
	Status() int
	Sent() uint64
	Recv() uint64
	Hijacked() bool
}

type responseLogger struct {
	w        http.ResponseWriter
	status   int
	sent     uint64
	recv     uint64
	hijacked bool
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	n, err := l.w.Write(b)
	if atomicOps {
		atomic.AddUint64(&l.sent, uint64(n))
	} else {
		l.sent += uint64(n)
	}
	return n, err
}

func (l *responseLogger) WriteHeader(status int) {
	l.w.WriteHeader(status)
	l.status = status
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Sent() uint64 {
	if atomicOps {
		return atomic.LoadUint64(&l.sent)
	}
	return l.sent
}

func (l *responseLogger) Recv() uint64 {
	if atomicOps {
		return atomic.LoadUint64(&l.recv)
	}
	return l.recv
}

func (l *responseLogger) Hijacked() bool {
	return l.hijacked
}

type hijackLogger struct {
	*responseLogger
}

func (l *hijackLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h := l.w.(http.Hijacker)
	conn, rw, err := h.Hijack()
	if err != nil {
		return conn, rw, err
	}
	if l.status == 0 {
		// The status will be StatusSwitchingProtocols if there was no error and WriteHeader has not been called yet
		l.status = http.StatusSwitchingProtocols
	}
	l.hijacked = true
	return &connLogger{conn, l.responseLogger}, rw, err
}

type connLogger struct {
	net.Conn
	l *responseLogger
}

func (l *connLogger) Write(b []byte) (int, error) {
	n, err := l.Conn.Write(b)
	if atomicOps {
		atomic.AddUint64(&l.l.sent, uint64(n))
	} else {
		l.l.sent += uint64(n)
	}
	return n, err
}

func (l *connLogger) Read(b []byte) (int, error) {
	n, err := l.Conn.Read(b)
	if atomicOps {
		atomic.AddUint64(&l.l.recv, uint64(n))
	} else {
		l.l.recv += uint64(n)
	}
	return n, err
}
