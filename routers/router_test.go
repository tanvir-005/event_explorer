package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"event_explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

func TestRegisteredAPIRoutes(t *testing.T) {
	previousApp := services.App
	services.Initialize("", "")
	t.Cleanup(func() {
		services.App = previousApp
	})

	tests := []struct {
		path       string
		wantStatus int
	}{
		{path: "/api/locations/autocomplete", wantStatus: http.StatusBadRequest},
		{path: "/api/locations/place-123", wantStatus: http.StatusBadRequest},
		{path: "/api/cache/invalidate", wantStatus: http.StatusBadRequest},
		{path: "/api/cache/invalidate-all", wantStatus: http.StatusOK},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("GET %s returned status %d, want %d", test.path, response.Code, test.wantStatus)
			}
		})
	}
}

func TestRegisteredCacheInvalidationScope(t *testing.T) {
	previousApp := services.App
	services.Initialize("", "")
	t.Cleanup(func() {
		services.App = previousApp
	})

	tests := []struct {
		name       string
		query      string
		wantMusic  bool
		wantSports bool
	}{
		{
			name:       "single category",
			query:      "city=New%20York&countryCode=US&category=Sports",
			wantMusic:  true,
			wantSports: false,
		},
		{
			name:       "whole city",
			query:      "city=New%20York&countryCode=US",
			wantMusic:  false,
			wantSports: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := services.App.Cache
			cache.Clear()
			cache.Set("New York:US:Music", nil)
			cache.Set("New York:US:Sports", nil)

			request := httptest.NewRequest(
				http.MethodGet,
				"/api/cache/invalidate?"+test.query,
				nil,
			)
			response := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("invalidation returned status %d, want %d", response.Code, http.StatusOK)
			}

			_, musicFound := cache.Get("New York:US:Music")
			_, sportsFound := cache.Get("New York:US:Sports")
			if musicFound != test.wantMusic || sportsFound != test.wantSports {
				t.Fatalf(
					"cache state after invalidation: Music=%t Sports=%t, want Music=%t Sports=%t",
					musicFound,
					sportsFound,
					test.wantMusic,
					test.wantSports,
				)
			}
		})
	}
}
