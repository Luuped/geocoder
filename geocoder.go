package geocoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultNominatimDomain = "nominatim.openstreetmap.org"
	defaultUserAgent       = "geopy/1.0"
)

var rejectedUserAgents = []string{
	"my-application",
	"my_app/1",
	"my_user_agent/1.0",
	"specify_your_app_name_here",
	defaultUserAgent,
}

// Geocoder represents a Nominatim geocoder with configuration options.
type Geocoder struct {
	Domain     string
	UserAgent  string
	Scheme     string
	Timeout    time.Duration
	Proxies    map[string]string
	API        string
	ReverseAPI string
}

// Location represents a geographical location with display name, latitude, and longitude.
type Location struct {
	DisplayName    string                 `json:"display_name"`
	Lat            string                 `json:"lat"`
	Lon            string                 `json:"lon"`
	AddressDetails map[string]interface{} `json:"address,omitempty"`
}

// NewGeocoder initializes a new Geocoder instance with the specified user agent.
func NewGeocoder(userAgent string) (*Geocoder, error) {
	if userAgent == "" {
		return nil, errors.New("user agent is required")
	}

	for _, ua := range rejectedUserAgents {
		if userAgent == ua {
			return nil, fmt.Errorf("using Nominatim with user agent %s is discouraged", userAgent)
		}
	}

	geocoder := &Geocoder{
		Domain:     defaultNominatimDomain,
		UserAgent:  userAgent,
		Scheme:     "https",
		Timeout:    10 * time.Second,
		Proxies:    make(map[string]string),
		API:        "https://" + defaultNominatimDomain + "/search",
		ReverseAPI: "https://" + defaultNominatimDomain + "/reverse",
	}

	return geocoder, nil
}

// constructURL constructs the full URL for the API request with the given base API and parameters.
func (g *Geocoder) constructURL(baseAPI string, params url.Values) string {
	return baseAPI + "?" + params.Encode()
}

// Geocode performs a geocoding request to the Nominatim API with the specified query parameters.
// It returns either a single Location or a slice of Locations based on the exactlyOne parameter.
func (g *Geocoder) Geocode(query map[string]string, exactlyOne bool) (interface{}, error) {
	params := url.Values{}
	for key, value := range query {
		params.Set(key, value)
	}
	params.Set("format", "jsonv2")
	if exactlyOne {
		params.Set("limit", "1")
	}

	url := g.constructURL(g.API, params)
	resp, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := json.Unmarshal(resp, &locations); err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, errors.New("no results found")
	}

	if exactlyOne {
		return &locations[0], nil
	}

	return locations, nil
}

// Reverse performs a reverse geocoding request to the Nominatim API with the specified latitude and longitude.
func (g *Geocoder) Reverse(lat, lon float64, exactlyOne bool) (*Location, error) {
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("lon", fmt.Sprintf("%f", lon))
	params.Set("format", "jsonv2")
	params.Set("addressdetails", "1")

	url := g.constructURL(g.ReverseAPI, params)
	resp, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var location Location
	if err := json.Unmarshal(resp, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

// makeRequest performs the HTTP GET request to the specified URL and returns the response body.
func (g *Geocoder) makeRequest(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: g.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", g.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API request failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
