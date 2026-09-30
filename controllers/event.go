package controllers

import (
	"event_explorer/models"
	"event_explorer/services"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type TicketProvider interface {
	GetEvent(eventID string) (*models.Event, error)
}

type EventController struct {
	beego.Controller
	Ticketmaster TicketProvider
}

func (c *EventController) List() {
	city := c.GetString("city")
	countryCode := c.GetString("countryCode")

	if city == "" || countryCode == "" {
		c.Data["Error"] = "Please select a city first."
		c.TplName = "listing.tpl"
		return
	}

	results := services.App.Events.GetEvents(city, countryCode)

	c.Data["City"] = city
	c.Data["CountryCode"] = countryCode
	c.Data["Results"] = results
	c.TplName = "listing.tpl"
}

func (c *EventController) Details() {
	eventID := c.Ctx.Input.Param(":eventId")

	if eventID == "" {
		c.Data["Error"] = "Invalid event."
		c.TplName = "details.tpl"
		return
	}

	provider := c.Ticketmaster

	if provider == nil {
		provider = services.App.Ticketmaster
	}

	event, err := provider.GetEvent(eventID)

	if err != nil {
		c.Data["Error"] = "Event not found."
		c.TplName = "details.tpl"
		return
	}

	c.Data["Event"] = event
	c.TplName = "details.tpl"
}

func validateTicketURL(rawURL string) error {
	ticketURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	if !strings.EqualFold(ticketURL.Scheme, "https") {
		return fmt.Errorf("ticket link must use HTTPS")
	}

	if !strings.EqualFold(ticketURL.Hostname(), "www.ticketmaster.com") {
		return fmt.Errorf("ticket provider is not approved")
	}

	return nil
}

func (c *EventController) Redirect() {
	eventID := c.Ctx.Input.Param(":eventId")

	if eventID == "" {
		c.Data["Error"] = "Invalid event."
		c.TplName = "unavailable.tpl"
		return
	}

	provider := c.Ticketmaster

	if provider == nil {
		provider = services.App.Ticketmaster
	}

	event, err := provider.GetEvent(eventID)

	if err != nil {
		c.Data["Error"] = "Unable to find the event."
		c.TplName = "unavailable.tpl"
		return
	}

	if event.TicketURL == "" {
		c.Data["Error"] = "Ticket link is unavailable."
		c.TplName = "unavailable.tpl"
		return
	}

	if err := validateTicketURL(event.TicketURL); err != nil {
		c.Data["Error"] = "Invalid ticket link."
		c.TplName = "unavailable.tpl"
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Location", event.TicketURL)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusFound)
}
