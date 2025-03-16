package geocoder

import (
	"fmt"
)

func ExampleGeocoder_Geocode() {
	geocoder, err := NewGeocoder("zip_code_locator")
	if err != nil {
		fmt.Println("Error creating geocoder:", err)
		return
	}

	// Example querying by postal code
	query := map[string]string{
		"postalcode": "90210",
		"country":    "US",
	}

	result, err := geocoder.Geocode(query, true)
	if err != nil {
		fmt.Println("Error geocoding:", err)
		return
	}

	location, ok := result.(*Location)
	if !ok {
		fmt.Println("Error: expected *Location type")
		return
	}

	fmt.Printf("Location found: %s\n", location.DisplayName)
	fmt.Printf("Coordinates: %s, %s\n", location.Lat, location.Lon)
}

func ExampleGeocoder_GeocodeMultiple() {
	geocoder, err := NewGeocoder("multi_location_finder")
	if err != nil {
		fmt.Println("Error creating geocoder:", err)
		return
	}

	// Example querying by city name
	query := map[string]string{
		"city":    "San Francisco",
		"country": "US",
	}

	result, err := geocoder.Geocode(query, false)
	if err != nil {
		fmt.Println("Error geocoding:", err)
		return
	}

	locations, ok := result.([]Location)
	if !ok {
		fmt.Println("Error: expected []Location type")
		return
	}

	fmt.Printf("Found %d locations\n", len(locations))
	for i, loc := range locations {
		fmt.Printf("%d. %s (%s, %s)\n", i+1, loc.DisplayName, loc.Lat, loc.Lon)
	}
}

func ExampleGeocoder_Reverse() {
	geocoder, err := NewGeocoder("reverse_geocoder")
	if err != nil {
		fmt.Println("Error creating geocoder:", err)
		return
	}

	// Example reverse geocoding (Beverly Hills coordinates)
	lat, lon := 34.0736, -118.4004

	location, err := geocoder.Reverse(lat, lon, true)
	if err != nil {
		fmt.Println("Error reverse geocoding:", err)
		return
	}

	fmt.Printf("Address: %s\n", location.DisplayName)
}
