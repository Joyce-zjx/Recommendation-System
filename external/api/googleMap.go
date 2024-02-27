package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	latLongURL  = "https://maps.googleapis.com/maps/api/geocode/json"
	distanceURL = "https://maps.googleapis.com/maps/api/distancematrix/json"
)

type MapUtilities interface {
	// GetLatLong(address string) (float64, float64, error)
	CalculateDistance(lat1, long1, lat2, long2 float64) (int, error)
}

func NewMapUtilities(apiKey string) MapUtilities {
	return &googleMapUtilities{
		apiKey: apiKey,
	}
}

type googleMapUtilities struct {
	apiKey string
}

// address format should be "street, city, state"
// Make a request to the Google Maps Geocoding API to get the latitude and longitude coordinates for the address
// func (g *googleMapUtilities) GetLatLong(address string) (float64, float64, error) {
// 	client := &http.Client{}
// 	encodedAddress := url.QueryEscape(address)
// 	response, err := client.Get(latLongURL + "?address=" + encodedAddress + "&key=" + g.apiKey)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	data, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	var geocodingResponse struct {
// 		Results []struct {
// 			Geometry struct {
// 				Location struct {
// 					Lat float64 `json:"lat"`
// 					Lng float64 `json:"lng"`
// 				} `json:"location"`
// 			} `json:"geometry"`
// 		} `json:"results"`
// 	}
// 	err = json.Unmarshal(data, &geocodingResponse)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	if len(geocodingResponse.Results) == 0 {
// 		return 0, 0, fmt.Errorf("no results found for address %s", address)
// 	}

// 	// Return the latitude and longitude coordinates for the address
// 	return geocodingResponse.Results[0].Geometry.Location.Lat, geocodingResponse.Results[0].Geometry.Location.Lng, nil
// }

// Make a request to the Google Maps Distance Matrix API to get the distance between two locations in meters
func (g *googleMapUtilities) CalculateDistance(lat1, long1, lat2, long2 float64) (int, error) {
	client := &http.Client{}
	origins := fmt.Sprintf("%f,%f", lat1, long1)
	destinations := fmt.Sprintf("%f,%f", lat2, long2)
	response, err := client.Get(distanceURL + "?origins=" + url.QueryEscape(origins) +
		"&destinations=" + url.QueryEscape(destinations) + "&key=" + g.apiKey)
	if err != nil {
		return 0, err
	}
	var distanceResponse struct {
		Rows []struct {
			Elements []struct {
				Distance struct {
					Value int `json:"value"`
				} `json:"distance"`
			}
		} `json:"rows"`
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(data, &distanceResponse)
	if err != nil {
		return 0, err
	}
	if len(distanceResponse.Rows) == 0 || len(distanceResponse.Rows[0].Elements) == 0 {
		return 0, fmt.Errorf("no results found for distance")
	}
	return distanceResponse.Rows[0].Elements[0].Distance.Value, nil
}
