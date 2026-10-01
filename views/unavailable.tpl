<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#fafbf8">
    <title>Tickets unavailable | Event Explorer</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body class="page-message">
    {{template "layouts/header.tpl" .}}

    <main id="main" class="container message-page">
        <p class="eyebrow"><span class="eyebrow-mark">A QUICK NOTE</span> TICKET INFORMATION</p>
        <span class="message-symbol" aria-hidden="true">!</span>
        <h1>Tickets unavailable</h1>
        <p>{{if .Error}}{{.Error}}{{else}}We couldn't open ticket information for this event right now.{{end}}</p>
        <a class="button button-primary" href="/">Back to discovery <span aria-hidden="true">&rarr;</span></a>
    </main>

    {{template "layouts/footer.tpl" .}}
</body>
</html>
