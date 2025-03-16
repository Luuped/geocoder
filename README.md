# Geocoder Package

A Go library for geocoding and reverse geocoding using the Nominatim OpenStreetMap API.

## Features

- Forward geocoding (address/location name to coordinates)
- Reverse geocoding (coordinates to address)
- Configurable timeout and user agent
- Multiple results support

## Installation

```bash
go get github.com/luuped/geocoder
```

## Usage

### Initializing the Geocoder

```go
package main

import (
	"fmt"
	"github.com/luuped/geocoder"
)

func main() {
	// Initialize with your application's user agent
	// IMPORTANT: Using a descriptive user agent is required by the Nominatim Usage Policy
	gc, err := geocoder.NewGeocoder("your-application-name")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	// Proceed with geocoding operations
}
```

### Forward Geocoding

#### Single Location

```go
// Find coordinates for a postal code
query := map[string]string{
	"postalcode": "90210",
	"country":    "US",
}

result, err := gc.Geocode(query, true) // true means return only one result
if err != nil {
	fmt.Println("Error geocoding:", err)
	return
}

location, ok := result.(*geocoder.Location)
if !ok {
	fmt.Println("Error: unexpected result type")
	return
}

fmt.Printf("Location: %s\n", location.DisplayName)
fmt.Printf("Coordinates: %s, %s\n", location.Lat, location.Lon)
```

#### Multiple Locations

```go
// Find all matches for a city name
query := map[string]string{
	"city":    "San Francisco",
	"country": "US",
}

result, err := gc.Geocode(query, false) // false means return multiple results
if err != nil {
	fmt.Println("Error geocoding:", err)
	return
}

locations, ok := result.([]geocoder.Location)
if !ok {
	fmt.Println("Error: unexpected result type")
	return
}

fmt.Printf("Found %d locations\n", len(locations))
for i, loc := range locations {
	fmt.Printf("%d. %s (%s, %s)\n", i+1, loc.DisplayName, loc.Lat, loc.Lon)
}
```

### Reverse Geocoding

```go
// Find address for coordinates
lat, lon := 34.0736, -118.4004  // Beverly Hills coordinates
location, err := gc.Reverse(lat, lon, true)
if err != nil {
	fmt.Println("Error reverse geocoding:", err)
	return
}

fmt.Printf("Address: %s\n", location.DisplayName)
// Access specific address components
fmt.Printf("Address details: %+v\n", location.AddressDetails)
```

## Configuration Options

The Geocoder type includes several configuration options:

```go
type Geocoder struct {
	Domain     string            // API domain (default: nominatim.openstreetmap.org)
	UserAgent  string            // Required user agent string
	Scheme     string            // HTTP scheme (default: https)
	Timeout    time.Duration     // Request timeout (default: 10 seconds)
	Proxies    map[string]string // Proxy settings (if needed)
	API        string            // API endpoint for geocoding
	ReverseAPI string            // API endpoint for reverse geocoding
}
```

## Important Notes

1. **User Agent Requirement**: Nominatim's usage policy requires a valid user agent that identifies your application. Using generic user agents will result in an error.

2. **Usage Policy**: Be sure to comply with the [Nominatim Usage Policy](https://operations.osmfoundation.org/policies/nominatim/):
   - No heavy usage (an absolute maximum of 1 request per second)
   - Provide a valid HTTP Referer or User-Agent identifying your application
   - Put a valid email address in the request if possible
   - Spread load over longer periods for bulk requests

3. **Error Handling**: Always check for errors when making requests.

## License

[MIT License](LICENSE) - See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.