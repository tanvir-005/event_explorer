package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const googlePlacesBaseURL = "https://places.googleapis.com"

type GooglePlacesService struct {
	apiKey string
	client  *http.Client
}

func NewGooglePlacesService(apiKey string) *GooglePlacesService {
	return &GooglePlacesService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type AutocompleteRequest struct {
	Input               string   `json:"input"`
	IncludedPrimaryTypes []string `json:"includedPrimaryTypes"`
	SessionToken        string   `json:"sessionToken"`
}

type AutocompleteResponse struct {
	Suggestions []struct {
		PlacePrediction struct {
			PlaceID string `json:"placeId"`
			Text    struct {
				Text string `json:"text"`
			} `json:"text"`
		} `json:"placePrediction"`
	} `json:"suggestions"`
}

type PlaceDetailsResponse struct {
	AddressComponents []struct {
		LongText  string   `json:"longText"`
		ShortText string   `json:"shortText"`
		Types     []string `json:"types"`
	} `json:"addressComponents"`
}

func (s *GooglePlacesService) Autocomplete(
	input string,
	sessionToken string,
) (*AutocompleteResponse, error) {

	body := AutocompleteRequest{
		Input:                input,
		IncludedPrimaryTypes: []string{"(cities)"},
		SessionToken:         sessionToken,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		googlePlacesBaseURL+"/v1/places:autocomplete",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google autocomplete returned status %d", resp.StatusCode)
	}

	var result AutocompleteResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *GooglePlacesService) GetPlace(
	placeID string,
	sessionToken string,
) (*PlaceDetailsResponse, error) {

	url := fmt.Sprintf(
		"%s/v1/places/%s?sessionToken=%s",
		googlePlacesBaseURL,
		placeID,
		sessionToken,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Goog-Api-Key", s.apiKey)
	req.Header.Set("X-Goog-FieldMask", "addressComponents")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google place details returned status %d", resp.StatusCode)
	}

	var result PlaceDetailsResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}