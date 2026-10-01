package services

import (
	"errors"
	"sync"
	"testing"

	"event_explorer/models"
)

type mockEventProvider struct {
	mu        sync.Mutex
	events    map[string][]models.Event
	err       error
	callCount int
}

func (m *mockEventProvider) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount++

	if m.err != nil {
		return nil, m.err
	}

	return m.events[category], nil
}

func (m *mockEventProvider) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.callCount
}

func TestEventServiceReturnsMusicAndSports(t *testing.T) {
	provider := &mockEventProvider{
		events: map[string][]models.Event{
			"Music": {
				{ID: "music-1", Name: "Music Event"},
			},
			"Sports": {
				{ID: "sports-1", Name: "Sports Event"},
			},
		},
	}

	cache := NewEventCache()
	service := NewEventService(provider, cache)

	results := service.GetEvents("Las Vegas", "US")

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	foundMusic := false
	foundSports := false

	for _, result := range results {
		if result.Category == "Music" {
			foundMusic = true

			if len(result.Events) != 1 {
				t.Fatalf("Music returned %d events, want 1", len(result.Events))
			}

			if result.Error != "" {
				t.Fatalf("Music returned unexpected error: %s", result.Error)
			}
		}

		if result.Category == "Sports" {
			foundSports = true

			if len(result.Events) != 1 {
				t.Fatalf("Sports returned %d events, want 1", len(result.Events))
			}

			if result.Error != "" {
				t.Fatalf("Sports returned unexpected error: %s", result.Error)
			}
		}
	}

	if !foundMusic {
		t.Fatal("Music result not found")
	}

	if !foundSports {
		t.Fatal("Sports result not found")
	}
}

func TestEventServiceCachesSuccessfulResults(t *testing.T) {
	provider := &mockEventProvider{
		events: map[string][]models.Event{
			"Music": {
				{ID: "music-1", Name: "Music Event"},
			},
			"Sports": {
				{ID: "sports-1", Name: "Sports Event"},
			},
		},
	}

	cache := NewEventCache()
	service := NewEventService(provider, cache)

	service.GetEvents("Las Vegas", "US")

	if provider.Calls() != 2 {
		t.Fatalf("first request made %d provider calls, want 2", provider.Calls())
	}

	service.GetEvents("Las Vegas", "US")

	if provider.Calls() != 2 {
		t.Fatalf(
			"second request made %d provider calls, want 2 because results should come from cache",
			provider.Calls(),
		)
	}
}

func TestEventServiceKeepsSuccessfulCategoryWhenOtherFails(t *testing.T) {
	provider := &selectiveMockEventProvider{
		events: map[string][]models.Event{
			"Music": {
				{ID: "music-1", Name: "Music Event"},
			},
		},
		errors: map[string]error{
			"Sports": errors.New("sports API failed"),
		},
	}

	cache := NewEventCache()
	service := NewEventService(provider, cache)

	results := service.GetEvents("Las Vegas", "US")

	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	for _, result := range results {
		switch result.Category {
		case "Music":
			if len(result.Events) != 1 {
				t.Fatalf("Music returned %d events, want 1", len(result.Events))
			}

			if result.Error != "" {
				t.Fatalf("Music returned unexpected error: %s", result.Error)
			}

		case "Sports":
			if result.Error == "" {
				t.Fatal("Sports should contain an error")
			}

			if len(result.Events) != 0 {
				t.Fatalf("Sports returned %d events, want 0", len(result.Events))
			}
		}
	}

	service.GetEvents("Las Vegas", "US")
	if provider.Calls() != 3 {
		t.Fatalf("provider made %d calls after retry, want 3 (cached Music and retried Sports)", provider.Calls())
	}
}

type selectiveMockEventProvider struct {
	mu        sync.Mutex
	events    map[string][]models.Event
	errors    map[string]error
	callCount int
}

func (m *selectiveMockEventProvider) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.callCount++

	if err, ok := m.errors[category]; ok {
		return nil, err
	}

	return m.events[category], nil
}

func (m *selectiveMockEventProvider) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.callCount
}

func TestInitializeWiresApplicationServices(t *testing.T) {
	previousApp := App
	t.Cleanup(func() {
		App = previousApp
	})

	Initialize("google-test-key", "ticketmaster-test-key")

	if App == nil || App.GooglePlaces == nil || App.Ticketmaster == nil || App.Events == nil || App.Cache == nil {
		t.Fatal("Initialize did not populate every application service")
	}
	if App.GooglePlaces.apiKey != "google-test-key" {
		t.Fatalf("got Google Places key %q, want google-test-key", App.GooglePlaces.apiKey)
	}
	if App.Ticketmaster.apiKey != "ticketmaster-test-key" {
		t.Fatalf("got Ticketmaster key %q, want ticketmaster-test-key", App.Ticketmaster.apiKey)
	}
	if App.Events.cache != App.Cache || App.Events.ticketmaster != App.Ticketmaster {
		t.Fatal("event service does not share the application cache and Ticketmaster service")
	}
}
