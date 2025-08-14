package main

import (
	"fmt"
	"net/http"
)

var authChan = make(chan SpotifyAuth)

func callBackHandler(w http.ResponseWriter, r *http.Request) {
	authChan <- SpotifyAuth{
		code:  r.URL.Query().Get("code"),
		state: r.URL.Query().Get("state"),
	}

	fmt.Fprintf(w, "Authentication successful! Window can be closed")
}
