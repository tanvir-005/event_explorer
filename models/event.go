package models

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Timezone    string `json:"timezone"`
	Venue       string `json:"venue"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Genre       string `json:"genre"`
	SalesStatus string `json:"salesStatus"`
	Description string `json:"description,omitempty"`
	PleaseNote  string `json:"pleaseNote,omitempty"`
	TicketLimit string `json:"ticketLimit,omitempty"`
	SeatmapURL  string `json:"seatmapUrl,omitempty"`
	TicketURL   string `json:"ticketUrl,omitempty"`
}

type EventResult struct {
	Category string  `json:"category"`
	Events   []Event `json:"events"`
	Error    string  `json:"error,omitempty"`
}
