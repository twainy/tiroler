package debug

import (
	"expvar"
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"syscall"

	"github.com/gorilla/mux"
)

const internal = false

var (
	profiling  bool
	goroutines = expvar.NewInt("goroutines")
)

func init() {
	flag.BoolVar(&profiling, "profiling", false, "enable profiling")
	expvar.Publish("blackbird/debug/pprof", expvar.Func(func() interface{} { return profiling }))
	go toggle()
}

// Register registers /debug handlers and returns the /debug Subrouter
func Register(r *mux.Router) *mux.Router {
	d := r.PathPrefix("/debug").MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
		return !internal || isInternal(r.RemoteAddr)
	}).Subrouter()
	// profiling
	p := d.PathPrefix("/pprof").MatcherFunc(func(*http.Request, *mux.RouteMatch) bool { return profiling }).Subrouter()
	p.Path("/profile").HandlerFunc(pprof.Profile)
	p.Path("/cmdline").HandlerFunc(pprof.Cmdline)
	p.Path("/symbol").HandlerFunc(pprof.Symbol)
	p.NewRoute().HandlerFunc(pprof.Index)
	// exported variables
	d.Path("/vars").HandlerFunc(expvarHandler)
	d.Path("/headers").HandlerFunc(headersHandler)
	return d
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	keys := make([]string, len(r.Header))
	i := 0
	for key := range r.Header {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	first := true
	fmt.Fprintf(w, "{")
	for _, key := range keys {
		if !first {
			fmt.Fprint(w, ",")
		}
		first = false
		fmt.Fprintf(w, "\n\t%q: %q", key, r.Header.Get(key))
	}
	fmt.Fprintf(w, "\n}\n")
}

func isInternal(addr string) bool {
	ip := net.ParseIP(addr)
	if ip.IsLoopback() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		// TODO: use megatron
		return ip4[0] == 10
	}
	return false
}

// toggle listens for SIGUSR1 signals and toggles pprof availability when received
func toggle() {
	c := notify(syscall.SIGUSR1)
	for {
		<-c
		profiling = !profiling
	}
}

func notify(s ...os.Signal) <-chan os.Signal {
	c := make(chan os.Signal, len(s))
	signal.Notify(c, s...)
	return c
}

func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	goroutines.Set(int64(runtime.NumGoroutine()))
	expvar.Do(func(kv expvar.KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}
