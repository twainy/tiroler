// Command example is a sample application built with Goji. Its goal is to give
// you a taste for what Goji looks like in the real world by artificially using
// all of its features.
//
// In particular, this is a complete working site for gritter.com, a site where
// users can post 140-character "greets". Any resemblance to real websites,
// alive or dead, is purely coincidental.
package tiroler

import (
//	"fmt"
//	"io"
	"net/http"
//	"regexp"
//	"strconv"
//	"time"
//	"encoding/json"

//	"github.com/goji/param"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
//	"github.com/zenazn/goji/web/middleware"
)

// Note: the code below cuts a lot of corners to make the example app simple.

func main() {
	// Add routes to the global handler
	goji.Get("/", Root)
	// Use Sinatra-style patterns in your URLs
	goji.Get("/novel/:ncode", GetNovel)

	// Middleware can be used to inject behavior into your app. The
	// middleware for this application are defined in middleware.go, but you
	// can put them wherever you like.
	goji.Use(Json)

	// Call Serve() at the bottom of your main() function, and it'll take
	// care of everything else for you, including binding to a socket (with
	// automatic support for systemd and Einhorn) and supporting graceful
	// shutdown on SIGINT. Serve() is appropriate for both development and
	// production.
	goji.Serve()
}

// Root route (GET "/"). Print a list of greets.
func Root(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(404), 404)
}

// GetUser finds a given user and her greets (GET "/user/:name")
func GetNovel(c web.C, w http.ResponseWriter, r *http.Request) {
	/*
	io.WriteString(w, "Gritter\n======\n\n")
	handle := c.URLParams["name"]
	user, ok := Users[handle]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	user.Write(w, handle)

	io.WriteString(w, "\nGreets:\n")
	for i := len(Greets) - 1; i >= 0; i-- {
		if Greets[i].User == handle {
			Greets[i].Write(w)
		}
	}
	*/
}

