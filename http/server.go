package http

import (
	"encoding/json"
	"fmt"
	"github.com/twainy/goban"
	"github.com/twainy/tiroler/crawler"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"net/http"
	"reflect"
	"runtime"
)

func Start() {
	Setup(goji.DefaultMux)
	goji.Serve()
}

func Setup(m *web.Mux) {
	route(m)
	goban.Setup("./conf/redis.json")
}

func route(m *web.Mux) {
	// Add routes to the global handler
	setGetHandler(m, "/", Root)
	// Use Sinatra-style patterns in your URLs
	setGetHandler(m, "/novel/:ncode", responseCache(getNovelInfo))

	// Middleware can be used to inject behavior into your app. The
	// middleware for this application are defined in middleware.go, but you
	// can put them wherever you like.
	m.Use(Json)
}

// Root route (GET "/"). Print a list of greets.
func Root(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

type NovelResponse struct {
	Tcode string
}

func setGetHandler(m *web.Mux, pattern interface{}, handler interface{}) {
	m.Get(pattern, handler)
}

func responseCache(handler func(c web.C, w http.ResponseWriter, r *http.Request) map[string]string) interface{} {
	return func(c web.C, w http.ResponseWriter, r *http.Request) {
		cache_key := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		for k, v := range c.URLParams {
			cache_key = cache_key + k + "_" + v
		}
		json_str, err := goban.Get(cache_key)
		if err != nil {
			response_map := handler(c, w, r)
			json_response, _ := json.Marshal(response_map)
			json_str = string(json_response)
			goban.Set(cache_key, string(json_str))
		}
		fmt.Fprint(w, string(json_str))
	}
}

// GetUser finds a given user and her greets (GET "/user/:name")
func getNovelInfo(c web.C, w http.ResponseWriter, r *http.Request) crawler.Novel {
	ncode := c.URLParams["ncode"]
	fmt.Println("get novel info novel", ncode)
	novel, _ := crawler.GetNovel(ncode)
    return novel
}

// PlainText sets the content-type of responses to text/plain.
func Json(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
