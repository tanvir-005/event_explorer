<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Events in {{.City}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    {{template "layouts/header.tpl" .}}

    <main>
        <h1>Events in {{.City}}, {{.CountryCode}}</h1>

        {{range .Results}}
            <section>
                <h2>{{.Category}}</h2>

                {{if .Error}}
                    <p>{{.Error}}</p>
                {{else if not .Events}}
                    <p>No {{.Category}} events found.</p>
                {{else}}
                    {{range .Events}}
                        <article>
                            {{if .ImageURL}}
                                <img
                                    src="{{.ImageURL}}"
                                    alt="{{.Name}}"
                                    width="300"
                                >
                            {{end}}

                            <h3>{{.Name}}</h3>
                            <p>
                                <strong>Date:</strong>
                                {{if .Date}}{{.Date}}{{else}}TBA{{end}}
                                |
                                {{if .Time}}{{.Time}}{{else}}TBA{{end}}
                            </p>
                            <p>
                                <strong>Venue:</strong>
                                {{if .Venue}}{{.Venue}}{{else}}TBA{{end}}
                                |
                                {{if .City}}{{.City}}{{else}}TBA{{end}}{{if .State}}, {{.State}}{{end}}
                            </p>

                            <a href="/events/{{.ID}}">
                                View Details
                            </a>
                        </article>
                    {{end}}
                {{end}}
            </section>
        {{end}}

        <a href="/">Search another city</a>
    </main>

    {{template "layouts/footer.tpl" .}}
</body>
</html>