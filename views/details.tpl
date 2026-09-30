<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .Event}}{{.Event.Name}}{{else}}Event Details{{end}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    {{template "layouts/header.tpl" .}}

    <main>
        <a href="javascript:history.back()">← Back</a>

        {{if .Error}}
            <h1>Event unavailable</h1>
            <p>{{.Error}}</p>
        {{else}}
            <article>
                {{if .Event.ImageURL}}
                    <img
                        src="{{.Event.ImageURL}}"
                        alt="{{.Event.Name}}"
                    >
                {{end}}

                <h1>{{.Event.Name}}</h1>

                <p>
                    <strong>Date:</strong>
                    {{if .Event.Date}}{{.Event.Date}}{{else}}TBA{{end}}
                    |
                    {{if .Event.Time}}{{.Event.Time}}{{else}}TBA{{end}}
                    {{if .Event.Timezone}}({{.Event.Timezone}}){{end}}
                </p>

                <p><strong>Venue:</strong> {{if .Event.Venue}}{{.Event.Venue}}{{else}}TBA{{end}}</p>
                {{if .Event.Address}}
                    <p>{{.Event.Address}}</p>
                {{end}}
                <p>
                    {{if .Event.City}}{{.Event.City}}{{else}}TBA{{end}}{{if .Event.State}}, {{.Event.State}}{{end}}
                </p>

                {{if .Event.Genre}}
                    <p><strong>Genre:</strong> {{.Event.Genre}}</p>
                {{end}}
                {{if .Event.SalesStatus}}
                    <p><strong>Sales status:</strong> {{.Event.SalesStatus}}</p>
                {{end}}

                {{if .Event.Description}}
                    <p>{{.Event.Description}}</p>
                {{end}}

                {{if .Event.PleaseNote}}
                    <section>
                        <h2>Event notes</h2>
                        <p>{{.Event.PleaseNote}}</p>
                    </section>
                {{end}}

                {{if .Event.TicketLimit}}
                    <p><strong>Ticket limit:</strong> {{.Event.TicketLimit}}</p>
                {{end}}

                {{if .Event.SeatmapURL}}
                    <section>
                        <h2>Seat map</h2>
                        <a href="{{.Event.SeatmapURL}}" target="_blank" rel="noopener noreferrer">
                            <img src="{{.Event.SeatmapURL}}" alt="Seat map for {{.Event.Name}}">
                        </a>
                    </section>
                {{end}}

                <a href="/redirect/{{.Event.ID}}">
                    View Tickets
                </a>
            </article>
        {{end}}
    </main>

    {{template "layouts/footer.tpl" .}}
</body>
</html>