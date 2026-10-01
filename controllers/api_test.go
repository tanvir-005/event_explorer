package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"event_explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type mockGooglePlacesProvider struct {
	autocompleteResult *services.AutocompleteResponse
	autocompleteError  error

	placeResult *services.PlaceDetailsResponse
	placeError  error

	autocompleteCalls int
	placeCalls        int
}

func (m *mockGooglePlacesProvider) Autocomplete(
	input string,
	sessionToken string,
) (*services.AutocompleteResponse, error) {
	m.autocompleteCalls++

	return m.autocompleteResult, m.autocompleteError
}

func (m *mockGooglePlacesProvider) GetPlace(
	placeID string,
	sessionToken string,
) (*services.PlaceDetailsResponse, error) {
	m.placeCalls++

	return m.placeResult, m.placeError
}

type mockCacheProvider struct {
	deletedKeys []string
	clearCalls  int
}

func (m *mockCacheProvider) Delete(key string) {
	m.deletedKeys = append(m.deletedKeys, key)
}

func (m *mockCacheProvider) Clear() {
	m.clearCalls++
}

func autocompleteResponse(t *testing.T) *services.AutocompleteResponse {
	t.Helper()

	data := `{
		"suggestions": [
			{
				"placePrediction": {
					"placeId": "place-dhaka",
					"text": {
						"text": "Dhaka, Bangladesh"
					}
				}
			},
			{
				"placePrediction": {
					"placeId": "",
					"text": {
						"text": "Invalid Place"
					}
				}
			},
			{
				"placePrediction": {
					"placeId": "place-chittagong",
					"text": {
						"text": "Chattogram, Bangladesh"
					}
				}
			}
		]
	}`

	var result services.AutocompleteResponse

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("failed to create autocomplete response: %v", err)
	}

	return &result
}

func placeDetailsResponse(t *testing.T) *services.PlaceDetailsResponse {
	t.Helper()

	data := `{
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
	}`

	var result services.PlaceDetailsResponse

	if err := json.Unmarshal([]byte(data), &result); err != nil {
		t.Fatalf("failed to create place details response: %v", err)
	}

	return &result
}

func newAPITestRegister(controller *APIController) *beego.ControllerRegister {
	register := beego.NewControllerRegister()

	register.Add(
		"/api/locations/autocomplete",
		controller,
		beego.WithRouterMethods(controller, "get:Autocomplete"),
	)

	register.Add(
		"/api/locations/:placeId",
		controller,
		beego.WithRouterMethods(controller, "get:Location"),
	)

	register.Add(
		"/api/cache/invalidate",
		controller,
		beego.WithRouterMethods(controller, "get:InvalidateCache"),
	)

	register.Add(
		"/api/cache/invalidate-all",
		controller,
		beego.WithRouterMethods(controller, "get:InvalidateAllCache"),
	)

	return register
}

func performAPIRequest(
	register *beego.ControllerRegister,
	method string,
	target string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	response := httptest.NewRecorder()

	register.ServeHTTP(response, request)

	return response
}

func TestAPIControllerAutocompleteSuccess(t *testing.T) {
	google := &mockGooglePlacesProvider{
		autocompleteResult: autocompleteResponse(t),
	}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/autocomplete?input=Dhaka&sessionToken=test-token",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var body struct {
		Suggestions []struct {
			PlaceID string `json:"placeId"`
			Text    string `json:"text"`
		} `json:"suggestions"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(body.Suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(body.Suggestions))
	}

	if body.Suggestions[0].PlaceID != "place-dhaka" {
		t.Errorf(
			"expected first place ID %q, got %q",
			"place-dhaka",
			body.Suggestions[0].PlaceID,
		)
	}

	if body.Suggestions[1].PlaceID != "place-chittagong" {
		t.Errorf(
			"expected second place ID %q, got %q",
			"place-chittagong",
			body.Suggestions[1].PlaceID,
		)
	}

	if google.autocompleteCalls != 1 {
		t.Errorf(
			"expected 1 autocomplete call, got %d",
			google.autocompleteCalls,
		)
	}
}

func TestAPIControllerAutocompleteMissingInput(t *testing.T) {
	google := &mockGooglePlacesProvider{}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/autocomplete?sessionToken=test-token",
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if google.autocompleteCalls != 0 {
		t.Errorf(
			"expected Google API not to be called, got %d calls",
			google.autocompleteCalls,
		)
	}
}

func TestAPIControllerAutocompleteMissingSessionToken(t *testing.T) {
	google := &mockGooglePlacesProvider{}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/autocomplete?input=Dhaka",
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if google.autocompleteCalls != 0 {
		t.Errorf(
			"expected Google API not to be called, got %d calls",
			google.autocompleteCalls,
		)
	}
}

func TestAPIControllerAutocompleteGoogleFailure(t *testing.T) {
	google := &mockGooglePlacesProvider{
		autocompleteError: errors.New("google API unavailable"),
	}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/autocomplete?input=Dhaka&sessionToken=test-token",
	)

	if response.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadGateway,
			response.Code,
		)
	}

	if google.autocompleteCalls != 1 {
		t.Errorf(
			"expected 1 autocomplete call, got %d",
			google.autocompleteCalls,
		)
	}
}

func TestAPIControllerLocationSuccess(t *testing.T) {
	google := &mockGooglePlacesProvider{
		placeResult: placeDetailsResponse(t),
	}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/place-dhaka?sessionToken=test-token",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	var body struct {
		City        string `json:"city"`
		CountryCode string `json:"countryCode"`
	}

	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.City != "Dhaka" {
		t.Errorf("expected city %q, got %q", "Dhaka", body.City)
	}

	if body.CountryCode != "BD" {
		t.Errorf(
			"expected country code %q, got %q",
			"BD",
			body.CountryCode,
		)
	}

	if google.placeCalls != 1 {
		t.Errorf(
			"expected 1 place lookup, got %d",
			google.placeCalls,
		)
	}
}

func TestAPIControllerLocationMissingSessionToken(t *testing.T) {
	google := &mockGooglePlacesProvider{}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/place-dhaka",
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if google.placeCalls != 0 {
		t.Errorf(
			"expected Google API not to be called, got %d calls",
			google.placeCalls,
		)
	}
}

func TestAPIControllerLocationGoogleFailure(t *testing.T) {
	google := &mockGooglePlacesProvider{
		placeError: errors.New("google API unavailable"),
	}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/place-dhaka?sessionToken=test-token",
	)

	if response.Code != http.StatusBadGateway {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadGateway,
			response.Code,
		)
	}

	if google.placeCalls != 1 {
		t.Errorf(
			"expected 1 place lookup, got %d",
			google.placeCalls,
		)
	}
}

func TestAPIControllerLocationMissingCityOrCountry(t *testing.T) {
	google := &mockGooglePlacesProvider{
		placeResult: &services.PlaceDetailsResponse{},
	}

	controller := &APIController{
		GooglePlaces: google,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/locations/place-dhaka?sessionToken=test-token",
	)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusUnprocessableEntity,
			response.Code,
		)
	}
}

func TestAPIControllerInvalidateCache(t *testing.T) {
	cache := &mockCacheProvider{}

	controller := &APIController{
		Cache: cache,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/cache/invalidate?city=Dhaka&countryCode=BD&category=Music",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if len(cache.deletedKeys) != 1 {
		t.Fatalf(
			"expected 1 deleted key, got %d",
			len(cache.deletedKeys),
		)
	}

	expectedKey := "Dhaka:BD:Music"

	if cache.deletedKeys[0] != expectedKey {
		t.Errorf(
			"expected deleted key %q, got %q",
			expectedKey,
			cache.deletedKeys[0],
		)
	}
}

func TestAPIControllerInvalidateCacheForCity(t *testing.T) {
	cache := &mockCacheProvider{}
	controller := &APIController{Cache: cache}
	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/cache/invalidate?city=Dhaka&countryCode=BD",
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}
	if len(cache.deletedKeys) != 2 {
		t.Fatalf("expected both category keys to be deleted, got %v", cache.deletedKeys)
	}
	wantKeys := []string{"Dhaka:BD:Music", "Dhaka:BD:Sports"}
	for index, wantKey := range wantKeys {
		if cache.deletedKeys[index] != wantKey {
			t.Errorf("deleted key %d = %q, want %q", index, cache.deletedKeys[index], wantKey)
		}
	}
}

func TestAPIControllerInvalidateCacheMissingParameters(t *testing.T) {
	cache := &mockCacheProvider{}

	controller := &APIController{
		Cache: cache,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/cache/invalidate?city=Dhaka&category=Music",
	)

	if response.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			response.Code,
		)
	}

	if len(cache.deletedKeys) != 0 {
		t.Errorf(
			"expected no cache deletion, got %d",
			len(cache.deletedKeys),
		)
	}
}

func TestAPIControllerInvalidateAllCache(t *testing.T) {
	cache := &mockCacheProvider{}

	controller := &APIController{
		Cache: cache,
	}

	register := newAPITestRegister(controller)

	response := performAPIRequest(
		register,
		http.MethodGet,
		"/api/cache/invalidate-all",
	)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}

	if cache.clearCalls != 1 {
		t.Errorf(
			"expected cache Clear to be called once, got %d",
			cache.clearCalls,
		)
	}
}

func TestAPIControllerResolvesDependenciesFromApp(t *testing.T) {
	previousApp := services.App
	services.Initialize("google-test-key", "ticketmaster-test-key")
	t.Cleanup(func() {
		services.App = previousApp
	})

	controller := &APIController{}
	if controller.googlePlaces() == nil {
		t.Fatal("expected Google Places dependency from the application container")
	}
	cache := controller.cache()
	if cache == nil {
		t.Fatal("expected cache dependency from the application container")
	}

	cache.Delete("test:US:Music")
	cache.Clear()
}
