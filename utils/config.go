package utils

import "os"

type Config struct {
	GooglePlacesAPIKey string
	TicketmasterAPIKey string
}

func LoadConfig() Config {
	return Config{
		GooglePlacesAPIKey: os.Getenv("GOOGLE_PLACES_API_KEY"),
		TicketmasterAPIKey: os.Getenv("TICKETMASTER_API_KEY"),
	}
}