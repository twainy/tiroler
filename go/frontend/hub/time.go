package hub

import (
	"time"

	"git.gree-dev.net/giistudio/blackbird-go/go/encoding/json"
	"github.com/golang/glog"
)

const (
	timePingType = "Time.Ping"
)

type ping struct {
	Ping int64 `json:"ping"`
}

type pong struct {
	Type string `json:"type"`
	Ping int64  `json:"ping"`
	Pong int64  `json:"pong"`
}

func timePing(c *conn, d []byte) error {
	var p ping
	err := json.Unmarshal(d, &p)
	if err != nil {
		glog.Infoln(c.LogTag(), "unmarshal", err)
		return err
	}
	d, err = json.Marshal(pong{
		Type: "Time.Pong",
		Ping: p.Ping,
		Pong: now(),
	})
	if err != nil {
		glog.Infoln(c.LogTag(), "marshal", err)
		return err
	}
	return c.Send(d)
}

func now() int64 {
	return time.Now().UnixNano() / 1000000
}
