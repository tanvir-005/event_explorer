package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"event_explorer/models"
)

const ticketmasterBaseURL = "https://app.ticketmaster.com/discovery/v2"

type TicketmasterService struct {
	apiKey string
	client *http.Client
}

func NewTicketmasterService(apiKey string) *TicketmasterService {
	return &TicketmasterService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ticketmasterResponse struct {
	Embedded struct {
		Events []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
			Dates struct {
				Start struct {
					LocalDate string `json:"localDate"`
				} `json:"localDate"`
			} `json:"dates"`
			Embedded struct {
				Venues []struct {
					Name string `json:"name"`
				} `json:"venues"`
			} `json:"_embedded"`
			URL string `json:"url"`
		} `json:"events"`
	} `json:"_embedded"`
}

func (s *TicketmasterService) GetEvents(
	city string,
	countryCode string,
	category string,
) ([]models.Event, error) {

	params := url.Values{}
	params.Set("apikey", s.apiKey)
	params.Set("city", city)
	params.Set("countryCode", countryCode)
	params.Set("classificationName", category)
	params.Set("size", "6")

	endpoint := ticketmasterBaseURL + "/events.json?" + params.Encode()

	resp, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ticketmaster returned status %d",
			resp.StatusCode,
		)
	}

	var result ticketmasterResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	events := make([]models.Event, 0)

	for _, item := range result.Embedded.Events {
		event := models.Event{
			ID:        item.ID,
			Name:      item.Name,
			Date:      item.Dates.Start.LocalDate,
			TicketURL: item.URL,
		}

		if len(item.Images) > 0 {
			event.ImageURL = item.Images[0].URL
		}

		if len(item.Embedded.Venues) > 0 {
			event.Venue = item.Embedded.Venues[0].Name
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *TicketmasterService) GetEvent(
	eventID string,
) (*models.Event, error) {
	params := url.Values{}
	params.Set("apikey", s.apiKey)

	endpoint := ticketmasterBaseURL +
		"/events/" + url.PathEscape(eventID) + ".json?" +
		params.Encode()

	resp, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ticketmaster returned status %d",
			resp.StatusCode,
		)
	}

	var item struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Images      []struct {
			URL string `json:"url"`
		} `json:"images"`
		Dates struct {
			Start struct {
				LocalDate string `json:"localDate"`
			} `json:"localDate"`
		} `json:"dates"`
		Embedded struct {
			Venues []struct {
				Name string `json:"name"`
			} `json:"venues"`
		} `json:"_embedded"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	event := &models.Event{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Date:        item.Dates.Start.LocalDate,
		TicketURL:   item.URL,
	}

	if len(item.Images) > 0 {
		event.ImageURL = item.Images[0].URL
	}

	if len(item.Embedded.Venues) > 0 {
		event.Venue = item.Embedded.Venues[0].Name
	}

	return event, nil
}
