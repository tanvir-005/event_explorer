package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"event_explorer/models"
	"event_explorer/services"

	beegoContext "github.com/beego/beego/v2/server/web/context"
)

type eventControllerTicketmaster struct {
	event       *models.Event
	err         error
	calls       int
	requestedID string
}

func (m *eventControllerTicketmaster) GetEvent(eventID string) (*models.Event, error) {
	m.calls++
	m.requestedID = eventID
	return m.event, m.err
}

type eventControllerProvider struct{}

func (eventControllerProvider) GetEvents(city, countryCode, category string) ([]models.Event, error) {
	return []models.Event{{ID: category + "-1", City: city}}, nil
}

func initEventControllerRequest(
	controller *EventController,
	method string,
	target string,
) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, target, nil)
	response := httptest.NewRecorder()
	ctx := beegoContext.NewContext()
	ctx.Reset(response, request)
	controller.Init(ctx, "EventController", "", nil)
	return response
}

func installEventService(t *testing.T) {
	t.Helper()
	previousApp := services.App
	services.App = &services.Container{
		Events: services.NewEventService(eventControllerProvider{}, services.NewEventCache()),
	}
	t.Cleanup(func() {
		services.App = previousApp
	})
}

func TestEventControllerListRequiresCityAndCountry(t *testing.T) {
	controller := &EventController{}
	initEventControllerRequest(controller, http.MethodGet, "/events?city=Dhaka")

	controller.List()

	if controller.TplName != "listing.tpl" {
		t.Fatalf("got template %q, want listing.tpl", controller.TplName)
	}
	if controller.Data["Error"] != "Please select a city first." {
		t.Fatalf("got error %v, want missing-location message", controller.Data["Error"])
	}
}

func TestEventControllerListLoadsResults(t *testing.T) {
	installEventService(t)
	controller := &EventController{}
	initEventControllerRequest(controller, http.MethodGet, "/events?city=Dhaka&countryCode=BD")

	controller.List()

	if controller.TplName != "listing.tpl" {
		t.Fatalf("got template %q, want listing.tpl", controller.TplName)
	}
	if controller.Data["City"] != "Dhaka" || controller.Data["CountryCode"] != "BD" {
		t.Fatalf("unexpected location data: city=%v countryCode=%v", controller.Data["City"], controller.Data["CountryCode"])
	}
	results, ok := controller.Data["Results"].([]models.EventResult)
	if !ok || len(results) != 2 {
		t.Fatalf("got results %#v, want two category results", controller.Data["Results"])
	}
}

func TestEventControllerDetailsRejectsMissingID(t *testing.T) {
	controller := &EventController{}
	initEventControllerRequest(controller, http.MethodGet, "/events/")

	controller.Details()

	if controller.TplName != "details.tpl" || controller.Data["Error"] != "Invalid event." {
		t.Fatalf("unexpected missing-ID response: template=%q data=%v", controller.TplName, controller.Data)
	}
}

func TestEventControllerDetailsReturnsProviderFailure(t *testing.T) {
	provider := &eventControllerTicketmaster{err: errors.New("provider unavailable")}
	controller := &EventController{Ticketmaster: provider}
	initEventControllerRequest(controller, http.MethodGet, "/events/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Details()

	if controller.TplName != "details.tpl" || controller.Data["Error"] != "Event not found." {
		t.Fatalf("unexpected provider-error response: template=%q data=%v", controller.TplName, controller.Data)
	}
	if provider.requestedID != "event-123" {
		t.Fatalf("provider received event ID %q, want event-123", provider.requestedID)
	}
}

func TestEventControllerDetailsReturnsEvent(t *testing.T) {
	wantEvent := &models.Event{ID: "event-123", Name: "Live show"}
	provider := &eventControllerTicketmaster{event: wantEvent}
	controller := &EventController{Ticketmaster: provider}
	initEventControllerRequest(controller, http.MethodGet, "/events/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Details()

	if controller.TplName != "details.tpl" || controller.Data["Event"] != wantEvent {
		t.Fatalf("unexpected successful response: template=%q event=%v", controller.TplName, controller.Data["Event"])
	}
}

func TestEventControllerRedirectRejectsMissingID(t *testing.T) {
	controller := &EventController{}
	initEventControllerRequest(controller, http.MethodGet, "/redirect/")

	controller.Redirect()

	if controller.TplName != "unavailable.tpl" || controller.Data["Error"] != "Invalid event." {
		t.Fatalf("unexpected missing-ID response: template=%q data=%v", controller.TplName, controller.Data)
	}
}

func TestEventControllerRedirectHandlesUnavailableEvent(t *testing.T) {
	provider := &eventControllerTicketmaster{err: errors.New("provider unavailable")}
	controller := &EventController{Ticketmaster: provider}
	initEventControllerRequest(controller, http.MethodGet, "/redirect/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Redirect()

	if controller.TplName != "unavailable.tpl" || controller.Data["Error"] != "Unable to find the event." {
		t.Fatalf("unexpected provider-error response: template=%q data=%v", controller.TplName, controller.Data)
	}
}

func TestEventControllerRedirectRejectsMissingTicketURL(t *testing.T) {
	provider := &eventControllerTicketmaster{event: &models.Event{ID: "event-123"}}
	controller := &EventController{Ticketmaster: provider}
	initEventControllerRequest(controller, http.MethodGet, "/redirect/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Redirect()

	if controller.Data["Error"] != "Ticket link is unavailable." {
		t.Fatalf("got error %v, want missing-ticket-link message", controller.Data["Error"])
	}
}

func TestEventControllerRedirectRejectsUnsafeTicketURL(t *testing.T) {
	provider := &eventControllerTicketmaster{event: &models.Event{TicketURL: "https://evil.example"}}
	controller := &EventController{Ticketmaster: provider}
	initEventControllerRequest(controller, http.MethodGet, "/redirect/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Redirect()

	if controller.Data["Error"] != "Invalid ticket link." {
		t.Fatalf("got error %v, want invalid-ticket-link message", controller.Data["Error"])
	}
}

func TestEventControllerRedirectsToValidTicketURL(t *testing.T) {
	ticketURL := "https://www.ticketmaster.com/event-123"
	provider := &eventControllerTicketmaster{event: &models.Event{TicketURL: ticketURL}}
	controller := &EventController{Ticketmaster: provider}
	response := initEventControllerRequest(controller, http.MethodGet, "/redirect/event-123")
	controller.Ctx.Input.SetParam(":eventId", "event-123")

	controller.Redirect()

	if response.Code != http.StatusFound {
		t.Fatalf("got status %d, want %d", response.Code, http.StatusFound)
	}
	if got := response.Header().Get("Location"); got != ticketURL {
		t.Fatalf("got redirect location %q, want %q", got, ticketURL)
	}
}

func TestHomeControllerSelectsHomeTemplate(t *testing.T) {
	controller := &HomeController{}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := beegoContext.NewContext()
	ctx.Reset(httptest.NewRecorder(), request)
	controller.Init(ctx, "HomeController", "Get", nil)

	controller.Get()

	if controller.TplName != "home.tpl" {
		t.Fatalf("got template %q, want home.tpl", controller.TplName)
	}
}

func TestValidateTicketURLAcceptsValidTicketmasterURL(t *testing.T) {
	err := validateTicketURL("https://www.ticketmaster.com/event-123")

	if err != nil {
		t.Fatalf("expected valid Ticketmaster URL, got error: %v", err)
	}
}

func TestValidateTicketURLRejectsHTTP(t *testing.T) {
	err := validateTicketURL("http://www.ticketmaster.com/event-123")

	if err == nil {
		t.Fatal("expected HTTP URL to be rejected")
	}
}

func TestValidateTicketURLRejectsUnapprovedHost(t *testing.T) {
	err := validateTicketURL("https://evil.example.com/event-123")

	if err == nil {
		t.Fatal("expected unapproved host to be rejected")
	}
}

func TestValidateTicketURLRejectsTicketmasterSubdomainAttack(t *testing.T) {
	err := validateTicketURL("https://www.ticketmaster.com.evil.example.com/event-123")

	if err == nil {
		t.Fatal("expected malicious hostname to be rejected")
	}
}

func TestValidateTicketURLRejectsMissingURL(t *testing.T) {
	err := validateTicketURL("")

	if err == nil {
		t.Fatal("expected empty URL to be rejected")
	}
}

func TestValidateTicketURLRejectsMalformedURL(t *testing.T) {
	err := validateTicketURL("://invalid-url")

	if err == nil {
		t.Fatal("expected malformed URL to be rejected")
	}
}
