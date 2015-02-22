package http

import (
	"net/http"
    "github.com/twainy/tiroler/api"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
    "encoding/json"
    "fmt"
)

func Start () {
    Route(goji.DefaultMux)
    goji.Serve()
}

func Route (m *web.Mux) {
	// Add routes to the global handler
	m.Get("/", Root)
	// Use Sinatra-style patterns in your URLs
	m.Get("/novel/:ncode", getNovelInfo)

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

// GetUser finds a given user and her greets (GET "/user/:name")
func getNovelInfo(c web.C, w http.ResponseWriter, r *http.Request) {
	ncode := c.URLParams["ncode"]
    fmt.Println("get novel info novel", ncode)
    tcode,err := api.GetTcode(ncode)
    response_map := map[string]string{"tcode":tcode}
    json_response,_ := json.Marshal(response_map)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
    fmt.Fprint(w, string(json_response))
}

// PlainText sets the content-type of responses to text/plain.
func Json(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}


