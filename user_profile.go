package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserProfile struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Href        string `json:"href"`
	ID          string `json:"id"`
	URI         string `json:"uri"`
}

func (cfg *ApiConfig) getUserID(token string) (string, error) {
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
