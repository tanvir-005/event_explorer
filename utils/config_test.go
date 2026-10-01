package utils

import "testing"

func TestLoadConfigReadsAPIKeysFromEnvironment(t *testing.T) {
	t.Setenv("GOOGLE_PLACES_API_KEY", "google-test-key")
	t.Setenv("TICKETMASTER_API_KEY", "ticketmaster-test-key")

	config := LoadConfig()

	if config.GooglePlacesAPIKey != "google-test-key" {
		t.Fatalf("got Google Places key %q, want google-test-key", config.GooglePlacesAPIKey)
	}
	if config.TicketmasterAPIKey != "ticketmaster-test-key" {
		t.Fatalf("got Ticketmaster key %q, want ticketmaster-test-key", config.TicketmasterAPIKey)
	}
}
