package config

import (
	"encoding/json"
	"expvar"
	"io/ioutil"
	"time"
)

var config = struct {
	ListenAddr string
	RPC        struct {
		URL                   string
		MaxIdleConnsPerHost   int
		DialTimeout           time.Duration
		ResponseHeaderTimeout time.Duration
		DisableKeepAlives     bool
		DisableCompression    bool
	}
}{}

func init() {
	expvar.Publish("tiroler/frontend/config", expvar.Func(func() interface{} { return config }))
}

func LoadJSON(data []byte) error {
	return json.Unmarshal(data, &config)
}

func LoadFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return LoadJSON(data)
}

func ListenAddr() string {
	return config.ListenAddr
}

func RpcURL() string {
	return config.RPC.URL
}

func RpcMaxIdleConnsPerHost() int {
	return config.RPC.MaxIdleConnsPerHost
}

func RpcDialTimeout() time.Duration {
	return config.RPC.DialTimeout
}

func RpcResponseHeaderTimeout() time.Duration {
	return config.RPC.ResponseHeaderTimeout
}

func RpcDisableKeepAlives() bool {
	return config.RPC.DisableKeepAlives
}

func RpcDisableCompression() bool {
	return config.RPC.DisableCompression
}
