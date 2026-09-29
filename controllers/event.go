package controllers

import (
	"event_explorer/services"

	beego "github.com/beego/beego/v2/server/web"
	"net/http"
	"net/url"
	"strings"
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
		c.CustomAbort(http.StatusBadRequest, "Invalid event.")
		return
	}

	event, err := services.App.Ticketmaster.GetEvent(eventID)
	if err != nil || event.TicketURL == "" {
		c.CustomAbort(http.StatusNotFound, "Ticket link not found.")
		return
	}

	parsedURL, err := url.Parse(event.TicketURL)
	if err != nil ||
		parsedURL.Scheme != "https" ||
		!strings.EqualFold(parsedURL.Hostname(), "www.ticketmaster.com") {
		c.CustomAbort(http.StatusBadGateway, "Invalid ticket link.")
		return
	}

	c.Ctx.ResponseWriter.Header().Set("Location", event.TicketURL)
	c.Ctx.ResponseWriter.WriteHeader(http.StatusFound)
}
