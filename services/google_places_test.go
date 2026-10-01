package services

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func googleServiceWithResponse(statusCode int, payload string) *GooglePlacesService {
	return &GooglePlacesService{
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

func TestAutocompleteReturnsResponse(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusOK,
		`{
			"suggestions": [
				{
					"placePrediction": {
						"placeId": "place-123",
						"text": {
							"text": "Dhaka, Bangladesh"
						}
					}
				}
			]
		}`,
	)

	response, err := service.Autocomplete("Dhaka", "session-123")

	if err != nil {
		t.Fatalf("Autocomplete returned error: %v", err)
	}

	if response == nil {
		t.Fatal("Autocomplete returned nil response")
	}

	if len(response.Suggestions) != 1 {
		t.Fatalf(
			"Autocomplete returned %d suggestions, want 1",
			len(response.Suggestions),
		)
	}

	suggestion := response.Suggestions[0]

	if suggestion.PlacePrediction.PlaceID != "place-123" {
		t.Fatalf(
			"unexpected place ID: %q",
			suggestion.PlacePrediction.PlaceID,
		)
	}

	if suggestion.PlacePrediction.Text.Text != "Dhaka, Bangladesh" {
		t.Fatalf(
			"unexpected suggestion text: %q",
			suggestion.PlacePrediction.Text.Text,
		)
	}
}

func TestGetPlaceReturnsResponse(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusOK,
		`{
			"addressComponents": [
				{
					"longText": "Dhaka",
					"shortText": "Dhaka",
					"types": ["locality"]
				},
				{
					"longText": "Bangladesh",
					"shortText": "BD",
					"types": ["country"]
				}
			]
		}`,
	)

	response, err := service.GetPlace("place-123", "session-123")

	if err != nil {
		t.Fatalf("GetPlace returned error: %v", err)
	}

	if response == nil {
		t.Fatal("GetPlace returned nil response")
	}

	if len(response.AddressComponents) != 2 {
		t.Fatalf(
			"GetPlace returned %d address components, want 2",
			len(response.AddressComponents),
		)
	}

	if response.AddressComponents[0].LongText != "Dhaka" {
		t.Fatalf(
			"unexpected city component: %q",
			response.AddressComponents[0].LongText,
		)
	}

	if response.AddressComponents[1].ShortText != "BD" {
		t.Fatalf(
			"unexpected country code: %q",
			response.AddressComponents[1].ShortText,
		)
	}
}

func TestAutocompleteReturnsErrorForAPIFailure(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusBadRequest,
		`{"error":{"message":"Invalid request"}}`,
	)

	response, err := service.Autocomplete("Dhaka", "session-123")

	if err == nil {
		t.Fatal("expected Autocomplete to return an error")
	}

	if response != nil {
		t.Fatalf(
			"expected nil response on API failure, got %+v",
			response,
		)
	}
}

func TestGetPlaceReturnsErrorForAPIFailure(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusNotFound,
		`{"error":{"message":"Place not found"}}`,
	)

	response, err := service.GetPlace("missing-place", "session-123")

	if err == nil {
		t.Fatal("expected GetPlace to return an error")
	}

	if response != nil {
		t.Fatalf(
			"expected nil response on API failure, got %+v",
			response,
		)
	}
}

func TestAutocompleteReturnsErrorForMalformedJSON(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusOK,
		`{"suggestions":`,
	)

	response, err := service.Autocomplete("Dhaka", "session-123")

	if err == nil {
		t.Fatal("expected Autocomplete to return an error for malformed JSON")
	}

	if response != nil {
		t.Fatalf(
			"expected nil response for malformed JSON, got %+v",
			response,
		)
	}
}

func TestGetPlaceReturnsErrorForMalformedJSON(t *testing.T) {
	service := googleServiceWithResponse(
		http.StatusOK,
		`{"addressComponents":`,
	)

	response, err := service.GetPlace("place-123", "session-123")

	if err == nil {
		t.Fatal("expected GetPlace to return an error for malformed JSON")
	}

	if response != nil {
		t.Fatalf(
			"expected nil response for malformed JSON, got %+v",
			response,
		)
	}
}

func TestNewGooglePlacesServiceConfiguresClient(t *testing.T) {
	service := NewGooglePlacesService("google-test-key")

	if service.apiKey != "google-test-key" {
		t.Fatalf("got API key %q, want google-test-key", service.apiKey)
	}
	if service.client == nil {
		t.Fatal("expected an HTTP client")
	}
	if service.client.Timeout != 10*time.Second {
		t.Fatalf("got client timeout %s, want 10s", service.client.Timeout)
	}
}

func TestGooglePlacesMethodsReturnTransportErrors(t *testing.T) {
	transportError := errors.New("network unavailable")
	service := &GooglePlacesService{
		client: &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, transportError
			}),
		},
	}

	if response, err := service.Autocomplete("Dhaka", "session-123"); err == nil || response != nil {
		t.Fatalf("Autocomplete returned response=%v, error=%v; want transport error", response, err)
	}
	if response, err := service.GetPlace("place-123", "session-123"); err == nil || response != nil {
		t.Fatalf("GetPlace returned response=%v, error=%v; want transport error", response, err)
	}
}

func TestGetPlaceEscapesPathAndSessionToken(t *testing.T) {
	var gotPath string
	var gotSessionToken string
	service := &GooglePlacesService{
		client: &http.Client{
			Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
				gotPath = request.URL.EscapedPath()
				gotSessionToken = request.URL.Query().Get("sessionToken")
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(`{"addressComponents":[]}`)),
					Request:    request,
				}, nil
			}),
		},
	}

	_, err := service.GetPlace("place/with space", "session&unexpected=yes")
	if err != nil {
		t.Fatalf("GetPlace returned error: %v", err)
	}
	if gotPath != "/v1/places/place%2Fwith%20space" {
		t.Fatalf("got escaped path %q, want escaped Place ID", gotPath)
	}
	if gotSessionToken != "session&unexpected=yes" {
		t.Fatalf("got session token %q, want the full token", gotSessionToken)
	}
}
