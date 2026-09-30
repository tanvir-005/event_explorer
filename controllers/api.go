package controllers

import (
	"event_explorer/models"
	"event_explorer/services"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"net/http"
	"slices"
	"strings"
)

type APIController struct {
	beego.Controller
}

func (c *APIController) respond(status int, data interface{}) {
	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = data
	c.ServeJSON()
}

func (c *APIController) Autocomplete() {
	input := strings.TrimSpace(c.GetString("input"))
	token := c.GetString("sessionToken")

	if input == "" || token == "" {
		c.respond(http.StatusBadRequest, map[string]string{
			"error": "input and sessionToken are required",
		})
		return
	}

	result, err := services.App.GooglePlaces.Autocomplete(input, token)

	if err != nil {
		c.respond(http.StatusBadGateway, map[string]string{
			"error": "Unable to fetch city suggestions",
		})
		return
	}

	suggestions := make([]models.CitySuggestion, 0)

	for _, item := range result.Suggestions {
		p := item.PlacePrediction

		if p.PlaceID == "" {
			continue
		}

		suggestions = append(suggestions, models.CitySuggestion{
			PlaceID: p.PlaceID,
			Text:    p.Text.Text,
		})
	}

	c.respond(http.StatusOK, map[string]interface{}{
		"suggestions": suggestions,
	})
}

func (c *APIController) Location() {
	placeID := c.Ctx.Input.Param(":placeId")
	token := c.GetString("sessionToken")

	if placeID == "" || token == "" {
		c.respond(http.StatusBadRequest, map[string]string{
			"error": "placeId and sessionToken are required",
		})
		return
	}

	result, err := services.App.GooglePlaces.GetPlace(placeID, token)

	if err != nil {
		c.respond(http.StatusBadGateway, map[string]string{
			"error": "Unable to resolve selected city",
		})
		return
	}

	location := models.Location{}

	for _, component := range result.AddressComponents {
		if slices.Contains(component.Types, "locality") {
			location.City = component.LongText
		}

		if slices.Contains(component.Types, "country") {
			location.CountryCode = component.ShortText
		}
	}

	if location.City == "" || location.CountryCode == "" {
		c.respond(http.StatusUnprocessableEntity, map[string]string{
			"error": "Unable to identify city or country",
		})
		return
	}

	c.respond(http.StatusOK, location)
}

func (c *APIController) InvalidateCache() {
	city := c.GetString("city")
	countryCode := c.GetString("countryCode")
	category := c.GetString("category")

	if city == "" || countryCode == "" || category == "" {
		c.respond(http.StatusBadRequest, map[string]string{
			"error": "city, countryCode and category are required",
		})
		return
	}

	key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)

	services.App.Cache.Delete(key)

	c.respond(http.StatusOK, map[string]string{
		"message": "cache invalidated",
	})
}

func (c *APIController) InvalidateAllCache() {
	services.App.Cache.Clear()

	c.respond(http.StatusOK, map[string]string{
		"message": "all cache data invalidated",
	})
}
