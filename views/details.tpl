<!DOCTYPE html>
<html>
<head>
    <title>{{.Event.Name}}</title>
</head>
<body>
    {{if .Error}}
        <h1>{{.Error}}</h1>
        <a href="/">Back</a>
    {{else}}
        <a href="javascript:history.back()">Back</a>

        <h1>{{.Event.Name}}</h1>

        {{if .Event.ImageURL}}
            <img src="{{.Event.ImageURL}}" alt="{{.Event.Name}}">
        {{end}}

        <p>Date: {{.Event.Date}}</p>
        <p>Venue: {{.Event.Venue}}</p>

        {{if .Event.Description}}
            <p>{{.Event.Description}}</p>
        {{end}}

        <a href="/redirect/{{.Event.ID}}">View Tickets</a>
    {{end}}
</body>
</html>