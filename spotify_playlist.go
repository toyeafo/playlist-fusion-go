package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
