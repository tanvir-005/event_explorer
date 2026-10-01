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
