package controllers

import (
	"event_explorer/models"
	"event_explorer/services"
	"fmt"
	"net/http"
	"slices"
	"strings"

	beego "github.com/beego/beego/v2/server/web"
)

type GooglePlacesProvider interface {
	Autocomplete(input, sessionToken string) (*services.AutocompleteResponse, error)
	GetPlace(placeID, sessionToken string) (*services.PlaceDetailsResponse, error)
}

type CacheProvider interface {
	Delete(key string)
	Clear()
}

type APIController struct {
	beego.Controller
	GooglePlaces GooglePlacesProvider
	Cache        CacheProvider
}

func (c *APIController) googlePlaces() GooglePlacesProvider {
	if c.GooglePlaces != nil {
		return c.GooglePlaces
	}

	return services.App.GooglePlaces
}

func (c *APIController) cache() CacheProvider {
	if c.Cache != nil {
		return c.Cache
	}

	return services.App.Cache
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

	result, err := c.googlePlaces().Autocomplete(input, token)

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

	result, err := c.googlePlaces().GetPlace(placeID, token)

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

	if city == "" || countryCode == "" {
		c.respond(http.StatusBadRequest, map[string]string{
			"error": "city and countryCode are required",
		})
		return
	}

	categories := []string{category}
	if category == "" {
		categories = []string{"Music", "Sports"}
	}

	for _, category := range categories {
		key := fmt.Sprintf("%s:%s:%s", city, countryCode, category)
		c.cache().Delete(key)
	}

	c.respond(http.StatusOK, map[string]string{
		"message": "cache invalidated",
	})
}

func (c *APIController) InvalidateAllCache() {
	c.cache().Clear()

	c.respond(http.StatusOK, map[string]string{
		"message": "all cache data invalidated",
	})
}
