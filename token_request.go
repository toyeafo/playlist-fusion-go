package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *ApiConfig) tokenRequest(authCode string) (SpotifyToken, error) {
	requestURL := "https://accounts.spotify.com/api/token"
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCode)
	formData.Set("redirect_uri", cfg.redirectURI)

	req, err := http.NewRequest(
		http.MethodPost,
		requestURL,
		bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return SpotifyToken{}, fmt.Errorf("error creating request, %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(cfg.spotifyClientID + ":" + cfg.spotifySecret))
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyToken{}, fmt.Errorf("error requesting token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SpotifyToken{}, fmt.Errorf("token request failed with status: %s", resp.Status)
	}

	var spotifyToken SpotifyToken
	if err := json.NewDecoder(resp.Body).Decode(&spotifyToken); err != nil {
		return SpotifyToken{}, fmt.Errorf("error decoding token response: %v", err)
	}

	return spotifyToken, nil
}
