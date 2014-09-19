package server

import (
	"errors"
	"flag"
	"fmt"
	"net/http"

	"github.com/twainy/tiroler-go/go/socket"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var configFilename string

func init() {
	flag.StringVar(&configFilename, "config", "", "configuration filename")
}

func ListenAndServe() error {
	if configFilename == "" {
		return errors.New("config: required")
	} else if err := config.LoadFile(configFilename); err != nil {
		return errors.New(fmt.Sprintf("config: %s", err))
	}

	r := mux.NewRouter()
	r.Path("/v1/ws").HandlerFunc(handleWs)
	r.Path("/health").HandlerFunc(handleHealth)
	debug.Register(r)
	glog.Infoln("listening on", config.ListenAddr())
	return http.ListenAndServe(config.ListenAddr(), logging.Handler(r))
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	ws := socket.NewWebSocket(w, r)
	if ws == nil {
		return
	}
	// @todo handling request
	<-ws.Closed() // block for logging
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
