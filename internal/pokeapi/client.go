package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the PokeAPI client
type Client struct {
	httpClient http.Client
	baseURL    string
}

// New creates a new PokeAPI client
func New() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		baseURL: "https://pokeapi.co/api/v2",
	}
}

// LocationAreaResponse represents the response from the location-area endpoint
type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// ListLocationAreas fetches location areas from the API
func (c *Client) ListLocationAreas(pageURL string) (LocationAreaResponse, error) {
	url := c.baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error fetching location areas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationAreaResponse{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var locationResp LocationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&locationResp)
	if err != nil {
		return LocationAreaResponse{}, fmt.Errorf("error parsing API response: %w", err)
	}

	return locationResp, nil
}
