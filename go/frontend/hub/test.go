package hub

import (
	"time"

	"git.gree-dev.net/giistudio/blackbird-go/go/encoding/json"
	"github.com/golang/glog"
)

const (
	testDelayType = "Test.Delay"
	testEchoType  = "Test.Echo"
)

type delay struct {
	Duration string `json:"duration"`
}

func testDelay(c *conn, d []byte) error {
	recv := time.Now()
	var delay delay
	if err := json.Unmarshal(d, &delay); err != nil {
		glog.Infoln(c.LogTag(), "unmarshal", err)
		return err
	}
	go func() {
		duration, _ := time.ParseDuration(delay.Duration)
		duration -= time.Since(recv)
		if duration > 0 {
			<-time.After(duration)
		}
		c.Send(d)
	}()
	return nil
}

func testEcho(c *conn, d []byte) error {
	return c.Send(d)
}
