package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type SpotifyAuth struct {
	code  string
	state string
}

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
	spotifyAuth     SpotifyAuth
	accessToken     string
	refreshToken    string
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
		grantType:       "authorization_code",
		redirectURI:     "http://127.0.0.1:8888/callback",
		scope:           "user-read-private user-read-email",
		responseType:    "code",
	}

	token, err := cfg.getAccessToken()
	if err != nil {
		log.Fatalf("authentication failed: %v", err)
	}

	fmt.Printf("Access Token: %s", token.AccessToken)
}
