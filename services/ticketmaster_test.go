package services

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
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

func serviceWithStatus(statusCode int, payload string) *TicketmasterService {
	return &TicketmasterService{
		apiKey: "test-key",
		client: &http.Client{
			Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: statusCode,
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

func TestGetEventsReturnsErrorForAPIFailure(t *testing.T) {
	service := serviceWithStatus(
		http.StatusUnauthorized,
		`{"fault":{"faultstring":"Invalid ApiKey"}}`,
	)

	events, err := service.GetEvents("Las Vegas", "US", "Music")

	if err == nil {
		t.Fatal("expected GetEvents to return an error")
	}

	if events != nil {
		t.Fatalf("expected no events on API failure, got %d", len(events))
	}
}

func TestGetEventReturnsErrorForAPIFailure(t *testing.T) {
	service := serviceWithStatus(
		http.StatusNotFound,
		`{"fault":{"faultstring":"Event not found"}}`,
	)

	event, err := service.GetEvent("missing-event")

	if err == nil {
		t.Fatal("expected GetEvent to return an error")
	}

	if event != nil {
		t.Fatalf("expected no event on API failure, got %+v", event)
	}
}

func TestNewTicketmasterServiceConfiguresClient(t *testing.T) {
	service := NewTicketmasterService("ticketmaster-test-key")

	if service.apiKey != "ticketmaster-test-key" {
		t.Fatalf("got API key %q, want ticketmaster-test-key", service.apiKey)
	}
	if service.client == nil {
		t.Fatal("expected an HTTP client")
	}
	if service.client.Timeout != 10*time.Second {
		t.Fatalf("got client timeout %s, want 10s", service.client.Timeout)
	}
}

func TestGetEventsReturnsErrorForTransportFailure(t *testing.T) {
	transportError := errors.New("network unavailable")
	service := &TicketmasterService{
		client: &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, transportError
			}),
		},
	}

	if events, err := service.GetEvents("Las Vegas", "US", "Music"); err == nil || events != nil {
		t.Fatalf("GetEvents returned events=%v, error=%v; want transport error", events, err)
	}
}

func TestGetEventsReturnsErrorForMalformedJSON(t *testing.T) {
	service := serviceWithResponse(`{"_embedded":`)

	if events, err := service.GetEvents("Las Vegas", "US", "Music"); err == nil || events != nil {
		t.Fatalf("GetEvents returned events=%v, error=%v; want decode error", events, err)
	}
}

func TestGetEventReturnsErrorForTransportFailure(t *testing.T) {
	transportError := errors.New("network unavailable")
	service := &TicketmasterService{
		client: &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, transportError
			}),
		},
	}

	if event, err := service.GetEvent("event-123"); err == nil || event != nil {
		t.Fatalf("GetEvent returned event=%v, error=%v; want transport error", event, err)
	}
}

func TestGetEventReturnsErrorForMalformedJSON(t *testing.T) {
	service := serviceWithResponse(`{"id":`)

	if event, err := service.GetEvent("event-123"); err == nil || event != nil {
		t.Fatalf("GetEvent returned event=%v, error=%v; want decode error", event, err)
	}
}
