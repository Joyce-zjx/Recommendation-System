package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pierrre/geohash"
)

const (
	discoveryURL             = "https://app.ticketmaster.com/discovery/v2/events.json"
	discoveryDEFAULT_KEYWORD = ""
	discoveryDEFAULT_RADIUS  = 500
	discoveryDEFAULT_UNIT    = "miles"
)

type EventsSearcher interface {
	ListEvents(latitude, longitude float64, radius int, keyword string) (TicketMasterDiscoveryResp, error)
}

type ticketMasterEventsSearcher struct {
	apiKey string
}

func NewEventsSearcher(apiKey string) EventsSearcher {
	return &ticketMasterEventsSearcher{
		apiKey: apiKey,
	}
}

type TicketMasterDiscoveryResp struct {
	Embedded embeddedEvents `json:"_embedded"`
}

type embeddedEvents struct {
	Events []event `json:"events"`
}
type event struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	URL             string           `json:"url"`
	Images          []image          `json:"images"`
	Dates           Dates            `json:"dates"`
	EmbeddedVenues  embeddedVenues   `json:"_embedded"`
	Classifications []classification `json:"classifications"`
}
type image struct {
	URL string `json:"url"`
}
type Dates struct {
	Start    start  `json:"start"`
	Timezone string `json:"timezone"`
	Status   status `json:"status"`
}
type status struct {
	Code string `json:"code"`
}
type start struct {
	LocalDate string `json:"localDate"`
	LocalTime string `json:"localTime"`
}
type classification struct {
	Genre    genre `json:"genre"`
	SubGenre genre `json:"subGenre"`
}
type genre struct {
	Name string `json:"name"`
}
type embeddedVenues struct {
	Venues []venue `json:"venues"`
}
type venue struct {
	Name     string   `json:"name"`
	Address  address  `json:"address"`
	City     city     `json:"city"`
	State    state    `json:"state"`
	Country  country  `json:"country"`
	Location location `json:"location"`
}
type city struct {
	Name string `json:"name"`
}
type state struct {
	Name string `json:"name"`
}
type country struct {
	Name string `json:"name"`
}
type address struct {
	Line1 string `json:"line1"`
}
type location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (t *ticketMasterEventsSearcher) ListEvents(latitude, longitude float64, radius int, keyword string) (TicketMasterDiscoveryResp, error) {
	if len(keyword) == 0 {
		keyword = discoveryDEFAULT_KEYWORD
	}
	if radius == 0 {
		radius = discoveryDEFAULT_RADIUS
	}
	client := &http.Client{}
	// function performs URL encoding in UTF-8, so we don't need to specify the encoding type explicitly.
	geoHash := geohash.Encode(latitude, longitude, 8)

	payload := url.Values{}
	payload.Add("apikey", t.apiKey)
	payload.Add("geoPoint", geoHash)
	payload.Add("keyword", keyword)
	payload.Add("unit", discoveryDEFAULT_UNIT)
	payload.Add("radius", fmt.Sprintf("%d", radius))
	req, err := http.NewRequest(http.MethodGet, discoveryURL+"?"+payload.Encode(), nil)
	if err != nil {
		return TicketMasterDiscoveryResp{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return TicketMasterDiscoveryResp{}, err
	}
	if res.StatusCode != http.StatusOK {
		return TicketMasterDiscoveryResp{}, errors.New("error status: " + res.Status)
	}
	ticketMasterRespBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return TicketMasterDiscoveryResp{}, err
	}
	defer res.Body.Close()
	var ticketMasterResp TicketMasterDiscoveryResp
	if err := json.Unmarshal(ticketMasterRespBytes, &ticketMasterResp); err != nil {
		return TicketMasterDiscoveryResp{}, err
	}

	return ticketMasterResp, nil
}
