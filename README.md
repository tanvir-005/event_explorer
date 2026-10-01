# Event Explorer

Event Explorer is a Fullstack web application developed in GO with Beego templates.

## Requirements And Configuration

Set these environment variables: get your api keys and put them in a `.env` file in the project root:

```dotenv
GOOGLE_PLACES_API_KEY=your_google_places_api_key
TICKETMASTER_API_KEY=your_ticketmaster_api_key
```

## Routes And API

All routes are registered in `routers/router.go`. All of the routes below use `GET`.

| Route | Parameters | Behavior |
| --- | --- | --- |
| `/` | None | Renders the home page. |
| `/events` | Query: `city`, `countryCode` | Renders Music and Sports event results for the location. Both parameters are required. |
| `/events/:eventId` | Path: `eventId` | Renders details for one Ticketmaster event. |
| `/redirect/:eventId` | Path: `eventId` | Looks up the event and redirects to its approved Ticketmaster ticket URL. |
| `/api/locations/autocomplete` | Query: `input`, `sessionToken` | Returns Google Places city suggestions as JSON. Both parameters are required. |
| `/api/locations/:placeId` | Path: `placeId`; query: `sessionToken` | Resolves a selected place to a JSON `city` and `countryCode`. |
| `/api/cache/invalidate` | Query: `city`, `countryCode`; optional `category` | Removes one category entry when `category` is supplied; otherwise removes both Music and Sports entries for that city/country. |
| `/api/cache/invalidate-all` | None | Removes every event-cache entry. |

## Event Cache

`services.Initialize` creates one `EventCache` when the application starts and shares it with `EventService`. The cache is in process memory and is protected by a `sync.RWMutex` so concurrent requests can read safely while updates and invalidations take an exclusive lock.

Events are cached separately for each city, country code, and category. The key format is:

```text
<city>:<countryCode>:<category>
```

For example, `Dhaka:BD:Music` and `Dhaka:BD:Sports` are different entries. On a miss, the service requests that category from Ticketmaster; the Music and Sports requests run concurrently. A successful response is cached, including an empty result. A failed request is returned for that category but not cached, so it will be retried on a later request. Cache hits avoid another Ticketmaster listing request.

There is no time-to-live or persistence: entries remain until the process exits, a matching entry is replaced, or an invalidation route is called. The cache is not shared across multiple application instances. Event detail lookups are made separately and are not stored in this event-list cache.

Invalidate one category with the exact location/category values used to form its key:

```sh
curl 'http://localhost:8080/api/cache/invalidate?city=Dhaka&countryCode=BD&category=Music'
```

Omit `category` to invalidate both Music and Sports entries for that city/country in one request:

```sh
curl 'http://localhost:8080/api/cache/invalidate?city=Dhaka&countryCode=BD'
```

Clear the complete cache for every city with:

```sh
curl 'http://localhost:8080/api/cache/invalidate-all'
```

## Tests

Command: `go test ./...`
```
ok      event_explorer                  0.010s
ok      event_explorer/controllers      (cached)
?       event_explorer/models           [no test files]
ok      event_explorer/routers          (cached)
ok      event_explorer/services         (cached)
ok      event_explorer/utils            (cached)
```
Command: `go test ./... -cover`
```
ok      event_explorer                  0.009s          coverage: 100.0% of statements
ok      event_explorer/controllers      (cached)        coverage: 98.3% of statements
?       event_explorer/models           [no test files]
ok      event_explorer/routers          (cached)        coverage: 100.0% of statements
ok      event_explorer/services         (cached)        coverage: 97.6% of statements
ok      event_explorer/utils            (cached)        coverage: 100.0% of statements
```