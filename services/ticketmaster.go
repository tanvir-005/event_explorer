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

type ticketmasterEvent struct {
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
			LocalTime string `json:"localTime"`
		} `json:"start"`
		Timezone string `json:"timezone"`
		Status   struct {
			Code string `json:"code"`
		} `json:"status"`
	} `json:"dates"`
	Classifications []struct {
		Primary bool `json:"primary"`
		Genre   struct {
			Name string `json:"name"`
		} `json:"genre"`
	} `json:"classifications"`
	PleaseNote  string `json:"pleaseNote"`
	TicketLimit struct {
		Info string `json:"info"`
	} `json:"ticketLimit"`
	Seatmap struct {
		StaticURL string `json:"staticUrl"`
	} `json:"seatmap"`
	Embedded struct {
		Venues []struct {
			Name    string `json:"name"`
			Address struct {
				Line1 string `json:"line1"`
			} `json:"address"`
			City struct {
				Name string `json:"name"`
			} `json:"city"`
			State struct {
				StateCode string `json:"stateCode"`
			} `json:"state"`
		} `json:"venues"`
	} `json:"_embedded"`
}

type ticketmasterResponse struct {
	Embedded struct {
		Events []ticketmasterEvent `json:"events"`
	} `json:"_embedded"`
}

func mapTicketmasterEvent(item ticketmasterEvent) models.Event {
	event := models.Event{
		ID:          item.ID,
		Name:        item.Name,
		Date:        item.Dates.Start.LocalDate,
		Time:        item.Dates.Start.LocalTime,
		Timezone:    item.Dates.Timezone,
		SalesStatus: item.Dates.Status.Code,
		Description: item.Description,
		PleaseNote:  item.PleaseNote,
		TicketLimit: item.TicketLimit.Info,
		SeatmapURL:  item.Seatmap.StaticURL,
		TicketURL:   item.URL,
	}

	if len(item.Images) > 0 {
		event.ImageURL = item.Images[0].URL
	}

	for _, classification := range item.Classifications {
		if classification.Primary || event.Genre == "" {
			event.Genre = classification.Genre.Name
		}
		if classification.Primary {
			break
		}
	}

	if len(item.Embedded.Venues) > 0 {
		venue := item.Embedded.Venues[0]
		event.Venue = venue.Name
		event.Address = venue.Address.Line1
		event.City = venue.City.Name
		event.State = venue.State.StateCode
	}

	return event
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
		events = append(events, mapTicketmasterEvent(item))
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

	var item ticketmasterEvent

	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, err
	}

	event := mapTicketmasterEvent(item)
	return &event, nil
}
