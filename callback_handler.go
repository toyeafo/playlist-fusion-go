package main

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) callBackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if state != cfg.spotifyAuth.state {
		http.Error(w, "State mismatch: Possible CSRF attack", http.StatusBadRequest)
		return
	}

	_, err := cfg.tokenRequest(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Authentication successful! Window can be closed")
}
