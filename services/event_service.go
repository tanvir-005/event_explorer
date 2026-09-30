package services

import (
	"fmt"

	"event_explorer/models"
)

type categoryResult struct {
	category string
	events   []models.Event
	err      error
}

type EventProvider interface {
	GetEvents(city string, countryCode string, category string) ([]models.Event, error)
}

type EventService struct {
	ticketmaster EventProvider
	cache        *EventCache
}

func NewEventService(
	ticketmaster EventProvider,
	cache *EventCache,
) *EventService {
	return &EventService{
		ticketmaster: ticketmaster,
		cache:        cache,
	}
}

func (s *EventService) GetEvents(
	city string,
	countryCode string,
) []models.EventResult {

	categories := []string{"Music", "Sports"}
	results := make(chan categoryResult, len(categories))

	for _, category := range categories {
		key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)

		if events, ok := s.cache.Get(key); ok {
			results <- categoryResult{
				category: category,
				events:   events,
			}
			continue
		}

		go func(category, key string) {
			events, err := s.ticketmaster.GetEvents(
				city,
				countryCode,
				category,
			)

			if err == nil {
				s.cache.Set(key, events)
			}

			results <- categoryResult{
				category: category,
				events:   events,
				err:      err,
			}
		}(category, key)
	}

	response := make([]models.EventResult, 0, len(categories))

	for i := 0; i < len(categories); i++ {
		result := <-results

		eventResult := models.EventResult{
			Category: result.category,
			Events:   result.events,
		}

		if result.err != nil {
			eventResult.Error = result.err.Error()
		}

		response = append(response, eventResult)
	}

	return response
}
