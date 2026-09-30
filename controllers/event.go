package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"event_explorer/services"

	beego "github.com/beego/beego/v2/server/web"
)

type EventController struct {
	beego.Controller
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

	event, err := services.App.Ticketmaster.GetEvent(eventID)
	if err != nil {
		c.Data["Error"] = "Event not found."
		c.TplName = "details.tpl"
		return
	}

	c.Data["Event"] = event
	c.TplName = "details.tpl"
}

func (c *EventController) Redirect() {
	eventID := c.Ctx.Input.Param(":eventId")

	if eventID == "" {
		c.Data["Error"] = "Invalid event."
		c.TplName = "unavailable.tpl"
		return
	}

	event, err := services.App.Ticketmaster.GetEvent(eventID)
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

	ticketURL, err := url.Parse(event.TicketURL)
	if err != nil {
		c.Data["Error"] = "Invalid ticket link."
		c.TplName = "unavailable.tpl"
		return
	}

	if ticketURL.Scheme != "https" {
		c.Data["Error"] = "Invalid ticket link."
		c.TplName = "unavailable.tpl"
		return
	}

	if !strings.EqualFold(ticketURL.Hostname(), "www.ticketmaster.com") {
		c.Data["Error"] = "Invalid ticket provider."
		c.TplName = "unavailable.tpl"
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Location", event.TicketURL)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusFound)
}
