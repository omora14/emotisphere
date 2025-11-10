package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LocationData struct {
	City    string
	Country string
	Lat     float64
	Lng     float64
}

type LocationService struct {
	Client *http.Client
}

func NewLocationService() *LocationService {
	return &LocationService{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetCoordinates gets coordinates for a city/country using Nominatim (OpenStreetMap)
// This is a free service that doesn't require an API key so I am good
func (ls *LocationService) GetCoordinates(city, country string) (float64, float64, error) {
	// Build query
	query := ""
	if city != "" {
		query = city
		if country != "" {
			query += ", " + country
		}
	} else if country != "" {
		query = country
	} else {
		return 0, 0, fmt.Errorf("both city and country are empty")
	}

	// Nominatim API (free, no API key required)
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1",
		strings.ReplaceAll(query, " ", "+"))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Emotisphere/1.0")

	resp, err := ls.Client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("Nominatim API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return 0, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no location found for: %s", query)
	}

	result := results[0]
	lat, ok1 := result["lat"].(string)
	lon, ok2 := result["lon"].(string)

	if !ok1 || !ok2 {
		return 0, 0, fmt.Errorf("invalid coordinate format in response")
	}

	var latFloat, lngFloat float64
	if _, err := fmt.Sscanf(lat, "%f", &latFloat); err != nil {
		return 0, 0, fmt.Errorf("failed to parse latitude: %w", err)
	}
	if _, err := fmt.Sscanf(lon, "%f", &lngFloat); err != nil {
		return 0, 0, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return latFloat, lngFloat, nil
}

func (ls *LocationService) ProcessLocation(countries []string) (string, string, error) {
	// extract first country if available
	if len(countries) > 0 && countries[0] != "" {
		countryName := mapCountryCode(countries[0])
		return "", countryName, nil
	}
	return "", "United States", nil
}

func mapCountryCode(code string) string {
	countryMap := map[string]string{
		"us": "United States",
		"cr": "Costa Rica",
		"br": "Brazil",
		"bo": "Bolivia",
		"es": "Spain",
		"gb": "United Kingdom",
		"jp": "Japan",
		"ca": "Canada",
		"au": "Australia",
		"de": "Germany",
		"fr": "France",
		"it": "Italy",
		"mx": "Mexico",
		"in": "India",
		"cn": "China",
		"ru": "Russia",
		"kr": "South Korea",
	}

	if name, ok := countryMap[strings.ToLower(code)]; ok {
		return name
	}
	return code
}
