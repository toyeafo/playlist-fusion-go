package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func (cfg *ApiConfig) getAccessToken() (SpotifyToken, error) {
	state, err := generateRandomState()
	if err != nil {
		return SpotifyToken{}, fmt.Errorf("failed to generate random state: %v", err)
	}

	authURL := userAuthRequest(cfg, state)
	browserAuth(authURL)

	fmt.Println("Please log into Spotify by visiting this URL in your browser:")
	fmt.Println(authURL)

	authChan := make(chan SpotifyAuth)

	server := &http.Server{Addr: "127.0.0.1:8000"}

	http.HandleFunc("/callback", callBackHandler)

	fmt.Println("Listening for a callback on http://127.0.0.1:8000/callback...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %s\n", err)
		}
	}()

	authResult := <-authChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	if authResult.state != state {
		log.Fatal("State mismatch.")
	}

	token, err := cfg.tokenRequest(authResult.code)
	if err != nil {
		return SpotifyToken{}, fmt.Errorf("error requesting token: %v", err)
	}

	return token, nil
}
