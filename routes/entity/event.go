package entity

// Event represents a business returned by the events searching API.
type Event struct {
	ID        string   `json:"item_id"`
	Name      string   `json:"name"`
	Date      Date     `json:"date"`
	Address   string   `json:"address"`
	Genres    []string `json:"genres"`
	ImageURL  string   `json:"image_url"`
	URL       string   `json:"url"`
	Longitude float64  `json:"longitude"`
	Latitude  float64  `json:"latitude"`
	Distance  int      `json:"distance"` // in meters
}

type Date struct {
	LocalDate string `json:"localDate"`
	LocalTime string `json:"localTime"`
	Timezone  string `json:"timezone"`
	Status    string `json:"status"`
}
