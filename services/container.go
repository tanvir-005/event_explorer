package services

type Container struct {
	GooglePlaces *GooglePlacesService
	Ticketmaster *TicketmasterService
	Events       *EventService
	Cache        *EventCache
}

var App *Container

func Initialize(googleKey, ticketmasterKey string) {
	cache := NewEventCache()
	ticketmaster := NewTicketmasterService(ticketmasterKey)

	App = &Container{
		GooglePlaces: NewGooglePlacesService(googleKey),
		Ticketmaster: ticketmaster,
		Cache:        cache,
		Events:       NewEventService(ticketmaster, cache),
	}
}
