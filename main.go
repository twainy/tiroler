// Command example is a sample application built with Goji. Its goal is to give
// you a taste for what Goji looks like in the real world by artificially using
// all of its features.
//
// In particular, this is a complete working site for gritter.com, a site where
// users can post 140-character "greets". Any resemblance to real websites,
// alive or dead, is purely coincidental.
package main

import (
    "github.com/twainy/tiroler/http"
)

// Note: the code below cuts a lot of corners to make the example app simple.

func main() {
	http.Start()
}
