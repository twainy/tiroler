package http

import (
    "testing"
    "github.com/zenazn/goji/web"
    "net/http/httptest"
    "net/http"
    "io/ioutil"
)

func ParseResponse(res *http.Response) (string, int) {
    defer res.Body.Close()
    contents, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }
    return string(contents), res.StatusCode
    
}

func TestGetNovelInfo(t *testing.T) {
    m := web.New()
    Route(m)
    ts := httptest.NewServer(m)
    defer ts.Close()
    
    res, err := http.Get(ts.URL + "/novel/n9902bs")
    if err != nil {
        t.Error("unexpected")
    }
    
    c, s := ParseResponse(res)
    if s != http.StatusOK {
        t.Error("invalid status code")
    }
    if c != `{"tcode":"449858"}` {
        t.Error("Invalid response.",  c)
    }
}