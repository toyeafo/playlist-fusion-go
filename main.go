package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const tokenRequestURL string = ""

type ApiConfig struct {
	spotifyClientID string
	spotifySecret   string
	geminiKey       string
	grantType       string
	scope           string
	state           string
	redirectURI     string
	responseType    string
	json            bool
}

type SpotifyToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

type SpotifyPlaylist struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
	Items    []struct {
		Collaborative bool   `json:"collaborative"`
		Description   string `json:"description"`
		ExternalUrls  struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href   string `json:"href"`
		ID     string `json:"id"`
		Images []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		Name  string `json:"name"`
		Owner struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Href        string `json:"href"`
			ID          string `json:"id"`
			Type        string `json:"type"`
			URI         string `json:"uri"`
			DisplayName string `json:"display_name"`
		} `json:"owner"`
		Public     bool   `json:"public"`
		SnapshotID string `json:"snapshot_id"`
		Tracks     struct {
			Href  string `json:"href"`
			Total int    `json:"total"`
		} `json:"tracks"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"items"`
}

type UserProfile struct {
	Country         string `json:"country"`
	DisplayName     string `json:"display_name"`
	Email           string `json:"email"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	URI     string `json:"uri"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error retrieving env variables")
	}

	spotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if spotifySecret == "" {
		log.Fatal("Spotify API secret must be set")
	}

	spotifyClientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if spotifyClientID == "" {
		log.Fatal("Spotify api client ID must be set")
	}

	geminiKey := os.Getenv("GEMINI_KEY")
	if geminiKey == "" {
		log.Fatal("Gemini key must be set")
	}

	cfg := &ApiConfig{
		spotifyClientID: spotifyClientID,
		spotifySecret:   spotifySecret,
		grantType:       "client_credentials",
	}

	result, err := cfg.tokenRequest()
	if err != nil {
		log.Fatalf("error requesting token: %v", err)
	}
	fmt.Println(result)
}

func (cfg *ApiConfig) tokenRequest() (SpotifyToken, error) {
	requestURL := "https://accounts.spotify.com/api/token"
	formData := url.Values{}
	formData.Set("grant_type", cfg.grantType)

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

	var spotifyToken SpotifyToken
	if err := json.NewDecoder(resp.Body).Decode(&spotifyToken); err != nil {
		return SpotifyToken{}, fmt.Errorf("error decoding token response: %v", err)
	}

	return spotifyToken, nil
}

func (cfg *ApiConfig) getUserProfile(token string) (string, error) {
	requestURL := "https://api.spotify.com/v1/me"

	req, err := http.NewRequest(
		http.MethodGet,
		requestURL,
		bytes.NewBufferString(""))
	if err != nil {
		return "", fmt.Errorf("error creating request, %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error retrieving user: %v", err)
	}
	defer resp.Body.Close()

	var userProfile UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		return "", fmt.Errorf("error decoding user profile: %v", err)
	}

	return userProfile.ID, nil
}

func (cfg *ApiConfig) getSpotifyPlaylist(userID, token string) (SpotifyPlaylist, error) {
	requestURL := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID)

	req, err := http.NewRequest(
		http.MethodGet,
		requestURL,
		bytes.NewBufferString(""))
	if err != nil {
		return SpotifyPlaylist{}, fmt.Errorf("error creating request, %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SpotifyPlaylist{}, fmt.Errorf("error requesting playlist: %v", err)
	}
	defer resp.Body.Close()

	var spotifyPlaylistResp SpotifyPlaylist
	if err := json.NewDecoder(resp.Body).Decode(&spotifyPlaylistResp); err != nil {
		return SpotifyPlaylist{}, fmt.Errorf("error decoding playlist response: %v", err)
	}

	return spotifyPlaylistResp, nil
}

func (cfg *ApiConfig) handleGeminiPlaylist() {}
