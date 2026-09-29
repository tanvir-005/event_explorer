<!DOCTYPE html>
<html>
<head>
    <title>Events</title>
</head>
<body>
    <h1>Events in {{.City}}</h1>

    {{if .Error}}
        <p>{{.Error}}</p>
    {{else}}
        {{range .Results}}
            <h2>{{.Category}}</h2>

            {{if .Error}}
                <p>{{.Error}}</p>
            {{else}}
                {{range .Events}}
                    <article>
                        <h3>{{.Name}}</h3>
                        <p>{{.Date}}</p>
                        <p>{{.Venue}}</p>
                    </article>
                {{end}}
            {{end}}
        {{end}}
    {{end}}
</body>
</html>