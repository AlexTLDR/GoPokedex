package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AlexTLDR/GoPokedex/internal/pokecache"
)

// Client is the PokeAPI client
type Client struct {
	httpClient http.Client
	baseURL    string
	cache      *pokecache.Cache
}

// New creates a new PokeAPI client
func New() Client {
	return Client{
		httpClient: http.Client{
			Timeout: time.Minute,
		},
		baseURL: "https://pokeapi.co/api/v2",
		cache:   pokecache.NewCache(5 * time.Minute),
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

type LocationArea struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// GetLocationArea fetches details about a specific location area
func (c *Client) GetLocationArea(name string) (LocationArea, error) {
	url := c.baseURL + "/location-area/" + name

	if cachedData, found := c.cache.Get(url); found {
		var locationArea LocationArea
		err := json.Unmarshal(cachedData, &locationArea)
		if err == nil {
			return locationArea, nil
		}
	}

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error fetching location area: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error reading response body: %w", err)
	}

	var locationArea LocationArea
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error parsing API response: %w", err)
	}

	c.cache.Add(url, body)

	return locationArea, nil
}

func (c *Client) ListLocationAreas(pageURL string) (LocationAreaResponse, error) {
	url := c.baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}

	if cachedData, found := c.cache.Get(url); found {
		var locationResp LocationAreaResponse
		err := json.Unmarshal(cachedData, &locationResp)
		if err == nil {
			return locationResp, nil
		}
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
