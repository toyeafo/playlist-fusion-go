package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const redirectURI string = "http://127.0.0.1:8888/callback"

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
		redirectURI:     redirectURI,
		scope:           "user-read-private user-read-email",
	}

	// result, err := cfg.tokenRequest()
	// if err != nil {
	// 	log.Fatalf("error requesting token: %v", err)
	// }
	// fmt.Println(result)
	http.HandleFunc("/callback", cfg.callBackHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
