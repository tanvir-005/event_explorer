package services

import (
	"log"
	"sync"

	"event_explorer/models"
)

type EventCache struct {
	mu      sync.RWMutex
	entries map[string][]models.Event
}

func NewEventCache() *EventCache {
	return &EventCache{
		entries: make(map[string][]models.Event),
	}
}

func (c *EventCache) Get(key string) ([]models.Event, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	events, ok := c.entries[key]

	if ok {
		log.Printf("cache hit: %s", key)
	}

	return events, ok
}

func (c *EventCache) Set(key string, events []models.Event) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = events
}

func (c *EventCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

func (c *EventCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string][]models.Event)
}
