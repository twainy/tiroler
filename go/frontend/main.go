package main

import (
    "flag"
	"runtime"

	"github.com/twainy/tiroler-go/frontend/server"
	"github.com/golang/glog"
)

func main() {
flag.Parse()
	glog.Infoln("set GOMAXPROCS from", runtime.GOMAXPROCS(runtime.NumCPU()), "to", runtime.GOMAXPROCS(0))
	glog.Fatalln(server.ListenAndServe())
}
