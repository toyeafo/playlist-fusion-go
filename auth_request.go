package main

import (
	"net/url"
)

func userAuthRequest(cfg *ApiConfig, state string) string {
	requestURL := "https://accounts.spotify.com/authorize?"
	formData := url.Values{}
	formData.Set("response_type", cfg.responseType)
	formData.Set("client_id", cfg.spotifyClientID)
	formData.Set("scope", cfg.scope)
	formData.Set("redirect_uri", cfg.redirectURI)
	formData.Set("state", state)

	return requestURL + formData.Encode()
}
