package main

import (
	"net/http"
	"net/url"
)

func (cfg *ApiConfig) userAuthRequest(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}
	cfg.spotifyAuth.state = state

	requestURL := "https://accounts.spotify.com/authorize?"
	formData := url.Values{}
	formData.Set("response_type", cfg.responseType)
	formData.Set("client_id", cfg.spotifyClientID)
	formData.Set("scope", cfg.scope)
	formData.Set("redirect_uri", cfg.redirectURI)
	formData.Set("state", state)

	authURL := requestURL + formData.Encode()
	http.Redirect(w, r, authURL, http.StatusFound)
}
