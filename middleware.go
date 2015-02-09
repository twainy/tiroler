package tiroler

import (
	"net/http"
)

// PlainText sets the content-type of responses to text/plain.
func Json(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

