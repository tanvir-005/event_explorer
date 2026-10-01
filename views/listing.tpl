<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#fafbf8">
    <title>Events in {{.City}} | Event Explorer</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body class="page-listing">
    {{template "layouts/header.tpl" .}}

    <main id="main" class="container listing-page">
        <div class="listing-intro">
            <div>
                <p class="eyebrow"><span class="eyebrow-mark">DISCOVER</span> EVENTS AROUND YOU</p>
                <h1>In <em>{{.City}}</em><span class="country-code">{{.CountryCode}}</span></h1>
                <p class="listing-subtitle">A little more live, right around the corner.</p>
            </div>
            <a class="button button-secondary" href="/">Search another city <span aria-hidden="true">&rarr;</span></a>
        </div>

        {{if .Error}}
            <section class="empty-state warning" role="alert">
                <span class="empty-symbol" aria-hidden="true">!</span>
                <h2>Choose a city to get started</h2>
                <p>{{.Error}}</p>
                <a class="text-link" href="/">Find a city <span aria-hidden="true">&rarr;</span></a>
            </section>
        {{else}}
            <nav class="listing-info" aria-label="Event categories">
                <span>Browse by category</span>
                {{range .Results}}
                    <a href="#{{.Category}}">{{.Category}} <span aria-hidden="true">&darr;</span></a>
                {{end}}
                <span class="results-location">{{.City}}, {{.CountryCode}}</span>
            </nav>

            {{range .Results}}
                <section class="event-section" id="{{.Category}}" aria-labelledby="heading-{{.Category}}">
                    <div class="section-heading">
                        <div>
                            <p class="mini-label">OUT IN THE CITY</p>
                            <h2 id="heading-{{.Category}}">{{.Category}}</h2>
                        </div>
                        {{if .Events}}<span class="result-count">{{len .Events}} events</span>{{end}}
                    </div>

                    {{if .Error}}
                        <div class="empty-state warning" role="status">
                            <span class="empty-symbol" aria-hidden="true">!</span>
                            <h3>{{.Category}} events aren't available right now</h3>
                            <p>{{.Error}}</p>
                        </div>
                    {{else if not .Events}}
                        <div class="empty-state">
                            <span class="empty-symbol" aria-hidden="true">+</span>
                            <h3>Nothing on the calendar just yet</h3>
                            <p>No {{.Category}} events found in {{$.City}}. Try another city to see what's happening.</p>
                            <a class="text-link" href="/">Search another city <span aria-hidden="true">&rarr;</span></a>
                        </div>
                    {{else}}
                        <div class="event-grid">
                            {{range .Events}}
                                <article class="event-card">
                                    <a class="card-image-link" href="/events/{{.ID}}" aria-label="View {{.Name}} details">
                                        {{if .ImageURL}}
                                            <img src="{{.ImageURL}}" alt="{{.Name}}" loading="lazy">
                                        {{else}}
                                            <span class="card-placeholder" aria-hidden="true"><span>LIVE / LOCAL</span></span>
                                        {{end}}
                                        {{if .Genre}}<span class="image-label">{{.Genre}}</span>{{end}}
                                    </a>
                                    <div class="card-body">
                                        <p class="event-date">{{if .Date}}{{.Date}}{{else}}Date to be announced{{end}}</p>
                                        <h3><a href="/events/{{.ID}}">{{.Name}}</a></h3>
                                        <p class="event-venue">{{if .Venue}}{{.Venue}}{{else}}Venue to be announced{{end}}</p>
                                        <div class="card-bottom">
                                            <span>{{if .Time}}{{.Time}}{{else}}Time TBA{{end}} <span aria-hidden="true">/</span> {{if .City}}{{.City}}{{else}}{{$.City}}{{end}}{{if .State}}, {{.State}}{{end}}</span>
                                            <a class="details-link" href="/events/{{.ID}}">Details <span aria-hidden="true">&rarr;</span></a>
                                        </div>
                                    </div>
                                </article>
                            {{end}}
                        </div>
                    {{end}}
                </section>
            {{end}}
            <p class="bottom-note">Event details are provided by Ticketmaster. Times are shown in the venue's local timezone when available.</p>
        {{end}}
    </main>

    {{template "layouts/footer.tpl" .}}
</body>
</html>