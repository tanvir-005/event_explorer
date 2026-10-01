<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#fafbf8">
    <title>{{if .Event}}{{.Event.Name}}{{else}}Event Details{{end}} | Event Explorer</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body class="page-details">
    {{template "layouts/header.tpl" .}}

    <main id="main" class="container detail-page">
        <a class="back-link" href="/">&larr; Back to discovery</a>

        {{if .Error}}
            <section class="empty-state warning detail-error" role="alert">
                <span class="empty-symbol" aria-hidden="true">!</span>
                <p class="mini-label">EVENT DETAILS</p>
                <h1>Event unavailable</h1>
                <p>{{.Error}}</p>
                <a class="button button-primary" href="/">Return to discovery <span aria-hidden="true">&rarr;</span></a>
            </section>
        {{else}}
            <article class="detail-grid">
                <div class="detail-main">
                    <div class="detail-image">
                        {{if .Event.ImageURL}}
                            <img src="{{.Event.ImageURL}}" alt="{{.Event.Name}}" fetchpriority="high">
                        {{else}}
                            <div class="detail-placeholder" aria-hidden="true"><span>LIVE / LOCAL</span></div>
                        {{end}}
                    </div>
                    {{if .Event.Genre}}<p class="detail-category"><span class="category-pill">{{.Event.Genre}}</span></p>{{end}}
                    <h1>{{.Event.Name}}</h1>

                    {{if .Event.Description}}
                        <section class="about-event">
                            <p class="mini-label">THE STORY</p>
                            <h2>About this event</h2>
                            <p class="description">{{.Event.Description}}</p>
                        </section>
                    {{end}}

                    {{if .Event.PleaseNote}}
                        <section class="detail-section">
                            <p class="mini-label">GOOD TO KNOW</p>
                            <h2>Event notes</h2>
                            <p>{{.Event.PleaseNote}}</p>
                        </section>
                    {{end}}

                    {{if .Event.SeatmapURL}}
                        <section class="detail-section seatmap-section">
                            <p class="mini-label">FIND YOUR VIEW</p>
                            <h2>Seat map</h2>
                            <a href="{{.Event.SeatmapURL}}" target="_blank" rel="noopener noreferrer">
                                <img src="{{.Event.SeatmapURL}}" alt="Seat map for {{.Event.Name}}" loading="lazy">
                                <span class="text-link">Open seat map <span aria-hidden="true">&nearr;</span></span>
                            </a>
                        </section>
                    {{end}}
                </div>

                <aside class="booking-panel" aria-label="Event information and tickets">
                    <p class="mini-label">MAKE A PLAN</p>
                    <h2>{{.Event.Name}}</h2>
                    <dl>
                        <div>
                            <dt>Date &amp; time</dt>
                            <dd>{{if .Event.Date}}{{.Event.Date}}{{else}}To be announced{{end}}<small>{{if .Event.Time}}{{.Event.Time}}{{else}}Time to be announced{{end}}{{if .Event.Timezone}} &middot; {{.Event.Timezone}}{{end}}</small></dd>
                        </div>
                        <div>
                            <dt>Venue</dt>
                            <dd>{{if .Event.Venue}}{{.Event.Venue}}{{else}}To be announced{{end}}
                                <small>{{if .Event.Address}}{{.Event.Address}}<br>{{end}}{{if .Event.City}}{{.Event.City}}{{else}}Location to be announced{{end}}{{if .Event.State}}, {{.Event.State}}{{end}}</small>
                            </dd>
                        </div>
                        {{if .Event.SalesStatus}}
                            <div><dt>Sales status</dt><dd><span class="sales-status">{{.Event.SalesStatus}}</span></dd></div>
                        {{end}}
                        {{if .Event.TicketLimit}}
                            <div><dt>Ticket limit</dt><dd>{{.Event.TicketLimit}}</dd></div>
                        {{end}}
                    </dl>
                    {{if .Event.TicketURL}}
                        <a class="button button-primary ticket-button" href="/redirect/{{.Event.ID}}">View tickets <span aria-hidden="true">&rarr;</span></a>
                        <p class="ticket-note">You'll continue to Ticketmaster to check availability and complete your purchase.</p>
                    {{else}}
                        <p class="unavailable-ticket">Ticket information isn't available yet.</p>
                    {{end}}
                    <p class="provider-line">EVENT INFORMATION BY TICKETMASTER</p>
                </aside>
            </article>
        {{end}}
    </main>

    {{template "layouts/footer.tpl" .}}
</body>
</html>