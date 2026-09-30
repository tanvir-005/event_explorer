package services

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (roundTrip roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

func serviceWithResponse(payload string) *TicketmasterService {
	return &TicketmasterService{
		apiKey: "test-key",
		client: &http.Client{
			Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(payload)),
					Request:    request,
				}, nil
			}),
		},
	}
}

func TestGetEventsMapsNestedDateAndVenue(t *testing.T) {
	service := serviceWithResponse(`{
		"_embedded": {
			"events": [{
				"id": "event-1",
				"name": "Live Show",
				"images": [{"url": "https://images.example/first.jpg"}],
				"dates": {
					"start": {"localDate": "2027-08-20", "localTime": "19:00:00"},
					"timezone": "America/Los_Angeles"
				},
				"_embedded": {"venues": [{
					"name": "Example Stadium",
					"city": {"name": "Las Vegas"},
					"state": {"stateCode": "NV"}
				}]}
			}]
		}
	}`)

	events, err := service.GetEvents("Las Vegas", "US", "Music")
	if err != nil {
		t.Fatalf("GetEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("GetEvents returned %d events, want 1", len(events))
	}

	event := events[0]
	if event.ID != "event-1" || event.ImageURL != "https://images.example/first.jpg" {
		t.Fatalf("event identity/image not mapped: %+v", event)
	}
	if event.Date != "2027-08-20" || event.Time != "19:00:00" || event.Timezone != "America/Los_Angeles" {
		t.Fatalf("nested event date not mapped: %+v", event)
	}
	if event.Venue != "Example Stadium" || event.City != "Las Vegas" || event.State != "NV" {
		t.Fatalf("venue location not mapped: %+v", event)
	}
}

func TestGetEventMapsDetailFieldsAndAllowsMissingOptionals(t *testing.T) {
	service := serviceWithResponse(`{
		"id": "event-1",
		"name": "Live Show",
		"url": "https://www.ticketmaster.com/event-1",
		"images": [{"url": "https://images.example/event.jpg"}],
		"dates": {
			"start": {"localDate": "2027-08-20", "localTime": "19:00:00"},
			"timezone": "America/Los_Angeles",
			"status": {"code": "onsale"}
		},
		"classifications": [{"primary": true, "genre": {"name": "Rock"}}],
		"pleaseNote": "Ticket transfers open later.",
		"ticketLimit": {"info": "Limit 4 per person."},
		"seatmap": {"staticUrl": "https://maps.example/seatmap.png"},
		"_embedded": {"venues": [{
			"name": "Example Stadium",
			"address": {"line1": "10 Main Street"},
			"city": {"name": "Las Vegas"},
			"state": {"stateCode": "NV"}
		}]}
	}`)

	event, err := service.GetEvent("event-1")
	if err != nil {
		t.Fatalf("GetEvent returned error: %v", err)
	}
	if event.Date != "2027-08-20" || event.Time != "19:00:00" || event.Timezone != "America/Los_Angeles" {
		t.Fatalf("nested event date not mapped: %+v", event)
	}
	if event.Address != "10 Main Street" || event.City != "Las Vegas" || event.State != "NV" {
		t.Fatalf("venue address/location not mapped: %+v", event)
	}
	if event.Genre != "Rock" || event.SalesStatus != "onsale" {
		t.Fatalf("classification or sales status not mapped: %+v", event)
	}
	if event.PleaseNote != "Ticket transfers open later." || event.TicketLimit != "Limit 4 per person." || event.SeatmapURL != "https://maps.example/seatmap.png" {
		t.Fatalf("ticket detail fields not mapped: %+v", event)
	}
	if event.Description != "" {
		t.Fatalf("missing description should stay empty, got %q", event.Description)
	}

	sparseService := serviceWithResponse(`{"id":"sparse-event","name":"Sparse Event"}`)
	sparseEvent, err := sparseService.GetEvent("sparse-event")
	if err != nil {
		t.Fatalf("GetEvent should accept missing optional fields: %v", err)
	}
	if sparseEvent.Date != "" || sparseEvent.Venue != "" || sparseEvent.SeatmapURL != "" {
		t.Fatalf("missing optional fields should remain empty: %+v", sparseEvent)
	}
}
