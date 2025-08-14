package main

import (
	"fmt"
	"log"
	"net/http"
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

	http.HandleFunc("/callback", callBackHandler)

	go func() {
		if err := http.ListenAndServe("127.0.0.1:8888", nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen failed: %s\n", err)
		}
	}()

	authResult := <-authChan

	if authResult.state != state {
		log.Fatal("State mismatch.")
	}

	token, err := cfg.tokenRequest(authResult.code)
	if err != nil {
		return SpotifyToken{}, fmt.Errorf("error requesting token: %v", err)
	}

	return token, nil
}
