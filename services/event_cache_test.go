package services

import (
	"testing"

	"event_explorer/models"
)

func TestEventCacheMiss(t *testing.T) {
	cache := NewEventCache()

	events, ok := cache.Get("Las Vegas:US:Music")

	if ok {
		t.Fatal("expected cache miss")
	}

	if events != nil {
		t.Fatalf("expected nil events on cache miss, got %+v", events)
	}
}

func TestEventCacheHit(t *testing.T) {
	cache := NewEventCache()

	expected := []models.Event{
		{ID: "event-1", Name: "Test Event"},
	}

	cache.Set("Las Vegas:US:Music", expected)

	events, ok := cache.Get("Las Vegas:US:Music")

	if !ok {
		t.Fatal("expected cache hit")
	}

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}

	if events[0].ID != "event-1" {
		t.Fatalf("got event ID %q, want %q", events[0].ID, "event-1")
	}
}

func TestEventCachePersistsUntilInvalidated(t *testing.T) {
	cache := NewEventCache()

	expected := []models.Event{
		{ID: "event-1", Name: "Test Event"},
	}

	key := "Las Vegas:US:Music"

	cache.Set(key, expected)

	events, ok := cache.Get(key)

	if !ok {
		t.Fatal("expected cached data to remain available")
	}

	if len(events) != 1 {
		t.Fatalf("got %d events, want 1", len(events))
	}
}

func TestEventCacheDelete(t *testing.T) {
	cache := NewEventCache()

	key := "Las Vegas:US:Music"

	cache.Set(key, []models.Event{
		{ID: "event-1"},
	})

	cache.Delete(key)

	_, ok := cache.Get(key)

	if ok {
		t.Fatal("expected cache miss after Delete")
	}
}

func TestEventCacheClear(t *testing.T) {
	cache := NewEventCache()

	cache.Set("Las Vegas:US:Music", []models.Event{
		{ID: "music-1"},
	})

	cache.Set("Las Vegas:US:Sports", []models.Event{
		{ID: "sports-1"},
	})

	cache.Clear()

	if _, ok := cache.Get("Las Vegas:US:Music"); ok {
		t.Fatal("Music cache entry still exists after Clear")
	}

	if _, ok := cache.Get("Las Vegas:US:Sports"); ok {
		t.Fatal("Sports cache entry still exists after Clear")
	}
}
