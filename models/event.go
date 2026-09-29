package models

type Event struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	Date        string `json:"date"`
	Venue       string `json:"venue"`
	Description string `json:"description,omitempty"`
	TicketURL   string `json:"ticketUrl,omitempty"`
}

type EventResult struct {
	Category string  `json:"category"`
	Events   []Event `json:"events"`
	Error    string  `json:"error,omitempty"`
}
