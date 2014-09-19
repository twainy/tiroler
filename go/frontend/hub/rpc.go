package hub

import (
	"bytes"
	j "encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"git.gree-dev.net/giistudio/blackbird-go/go/encoding/json"
	"git.gree-dev.net/giistudio/blackbird-go/go/frontend/config"
	"git.gree-dev.net/giistudio/blackbird-go/go/frontend/presence"
	"github.com/golang/glog"
)

const (
	rpcRequestType  = "RPC.Request"
	rpcResponseType = "RPC.Response"
	authMethod      = "Auth.Device"
)

var (
	errInternalService = j.RawMessage([]byte(`{"code":100001,"message":"INTERNAL_SERVER_ERROR"}`))
	errNotAllowed      = j.RawMessage([]byte(`{"code":100401,"message":"NOT_ALLOWED"}`))
	errUnderAttack     = j.RawMessage([]byte(`{"code":100101,"message":"UNDER_ATTACK"}`))
)

var httpClient = http.Client{
	Transport: &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, config.RpcDialTimeout())
		},
		DisableKeepAlives:     config.RpcDisableKeepAlives(),
		DisableCompression:    config.RpcDisableCompression(),
		MaxIdleConnsPerHost:   config.RpcMaxIdleConnsPerHost(),
		ResponseHeaderTimeout: config.RpcResponseHeaderTimeout(),
	},
}

type rpcClientRequest struct {
	ID     uint32        `json:"id"`
	Method string        `json:"method"`
	Args   *j.RawMessage `json:"args"`

	seq    uint32
	read   bool
	flight bool
}

type rpcClientResponse struct {
	Type   string        `json:"type"`
	ID     uint32        `json:"id"`
	Result *j.RawMessage `json:"result,omitempty"`
	Error  *j.RawMessage `json:"error,omitempty"`
}

type rpcServerRequest struct {
	Method  string        `json:"method"`
	Args    *j.RawMessage `json:"args"`
	Session *j.RawMessage `json:"session"`
}

type rpcServerResponse struct {
	AuthID   uint64        `json:"player_id"`
	Channels []string      `json:"channels"`
	Session  *j.RawMessage `json:"session"`
	Result   *j.RawMessage `json:"result"`
	Error    *j.RawMessage `json:"error"`
}

func rpcRequest(c *conn, d []byte) error {
	var r rpcClientRequest
	r.seq = c.seq.Add(1)
	if err := json.Unmarshal(d, &r); err != nil {
		return err
	}
	if m := strings.SplitN(r.Method, ".", 2); len(m) == 2 && len(m[1]) > 2 && m[1][:3] == "Get" {
		r.read = true
	}

	c.rpcLock.Lock()
	defer c.rpcLock.Unlock()
	if !allowed(c, r.Method) {
		return rpcResponse(c, &rpcClientResponse{
			ID:    r.ID,
			Error: &errNotAllowed,
		})
	}
	c.requests = append(c.requests, &r)
	sendRequests(c)
	return nil
}

func allowed(c *conn, method string) bool {
	if method == authMethod {
		return !authenticated(c) && len(c.requests) == 0
	}
	return authenticated(c)
}

func authenticated(c *conn) bool {
	return c.authID.Get() != 0
}

func rpcResponse(c *conn, r *rpcClientResponse) error {
	r.Type = rpcResponseType
	d, err := json.Marshal(r)
	if err != nil {
		glog.Errorln(c.LogTag(), "marshal", err)
		return err
	}
	return c.Send(d)
}

func sendRequests(c *conn) {
	for _, r := range c.requests {
		if !r.flight && canSend(c, r) {
			sendRequest(c, r)
		}
	}
}

func sendRequest(c *conn, r *rpcClientRequest) {
	r.flight = true
	// marshal the body with the session in the calling (locked) goroutine
	reqBody, err := json.Marshal(rpcServerRequest{
		Method:  r.Method,
		Args:    r.Args,
		Session: c.session,
	})
	// handle sending the request in a new goroutine
	go func() {
		var sResp rpcServerResponse
		defer func() {
			c.rpcLock.Lock()
			defer c.rpcLock.Unlock()
			if err := recover(); err != nil {
				glog.Errorln(c.LogTag(), "send", err)
				sResp.Result = nil
				sResp.Error = &errInternalService
			}
			if rpcResponse(c, &rpcClientResponse{
				ID:     r.ID,
				Result: sResp.Result,
				Error:  sResp.Error,
			}) != nil {
				c.Close(nil)
				// TODO: make sure client is disconnected from hub
				return
			}
			for i, req := range c.requests {
				if req.seq == r.seq {
					copy(c.requests[i:], c.requests[i+1:])                        // shift
					c.requests[len(c.requests)-1] = nil                           // clear
					c.requests = c.requests[:len(c.requests)-1 : cap(c.requests)] // reslice
					break
				}
			}
			if sResp.AuthID != 0 && c.authID.Get() == 0 {
				p := presence.NewPlayer(sResp.AuthID)
				if err := p.Login(c.Closed()); err != nil {
					// not technically correct... could be some problem with dynamodb, but this is sufficient for now
					c.Close(errUnderAttack)
					// TODO: make sure client is disconnected from hub
					return
				}
				c.authID.Set(sResp.AuthID)
			}
			if sResp.Session != nil {
				c.session = sResp.Session
			}
			sendRequests(c)
		}()
		if err != nil {
			// should never happen
			panic(err)
		}
		httpReq, err := http.NewRequest("POST", config.RpcURL(), bytes.NewBuffer(reqBody))
		if err != nil {
			// should never happen
			panic(err)
		}
		httpReq.Header.Add("Content-Type", "application/json")
		httpReq.Header.Add("X-Forwarded-For", c.RemoteAddr())
		// todo: authentication header?
		httpResp, err := httpClient.Do(httpReq)
		if err != nil {
			panic(err)
		}
		respBody, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			panic(err)
		}
		if httpResp.StatusCode != http.StatusOK {
			panic(httpResp.Status)
		}
		if err = json.Unmarshal(respBody, &sResp); err != nil {
			panic(err)
		}
	}()
}

func canSend(c *conn, r *rpcClientRequest) bool {
	if !authenticated(c) {
		return r.Method == authMethod
	}
	if r.read {
		return true
	}
	for _, r := range c.requests {
		if r.flight && !r.read {
			return false
		}
	}
	return true
}
