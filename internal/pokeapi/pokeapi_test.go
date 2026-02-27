package pokeapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetPokeLocations_Success(t *testing.T) {
	// fake API response
	mockResponse := Location{
		Count:    2,
		Next:     ptr("next-url"),
		Previous: nil,
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "canalave-city-area", URL: "url1"},
			{Name: "eterna-city-area", URL: "url2"},
		},
	}

	data, err := json.Marshal(mockResponse)
	if err != nil {
		t.Fatalf("failed to marshal mock response: %v", err)
	}

	// test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}))
	defer server.Close()

	client := NewClient(5*time.Second, 5*time.Minute)

	locations, urls, err := client.GetPokeLocations(server.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// assertions
	if len(locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locations))
	}

	if locations[0] != "canalave-city-area" {
		t.Errorf("unexpected first location: %s", locations[0])
	}

	if urls[0] != "next-url" {
		t.Errorf("expected next url 'next-url', got %s", urls[0])
	}

	if urls[1] != "" {
		t.Errorf("expected empty previous url, got %s", urls[1])
	}
}
func TestGetPokeLocations_UsesCache(t *testing.T) {
	client := NewClient(5*time.Second, 5*time.Minute)

	mockResponse := Location{
		Count: 1,
		Results: []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}{
			{Name: "cached-location", URL: "cached-url"},
		},
	}

	data, _ := json.Marshal(mockResponse)

	testURL := "https://cached.test"
	client.cache.Add(testURL, data)

	locations, _, err := client.GetPokeLocations(testURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(locations) != 1 || locations[0] != "cached-location" {
		t.Fatalf("cache not used properly, got %+v", locations)
	}
}
func ptr(s string) *string {
	return &s
}
