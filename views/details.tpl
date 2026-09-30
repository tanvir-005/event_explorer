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

                {{if .Event.Date}}
                    <p><strong>Date:</strong> {{.Event.Date}}</p>
                {{end}}

                {{if .Event.Venue}}
                    <p><strong>Venue:</strong> {{.Event.Venue}}</p>
                {{end}}

                {{if .Event.Description}}
                    <p>{{.Event.Description}}</p>
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