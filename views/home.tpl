<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Explorer</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    {{template "layouts/header.tpl" .}}

    <main>
        <h1>Event Explorer</h1>

        <form id="city-form">
            <input
                type="text"
                id="city-input"
                placeholder="Search for a city"
                autocomplete="off"
            >

            <div id="suggestions"></div>

            <input type="hidden" id="place-id">
            <input type="hidden" id="session-token">

            <button type="submit">Search</button>
        </form>
        <p class="google-attribution">Powered by Google</p>

        <p id="error"></p>
    </main>

    {{template "layouts/footer.tpl" .}}

    <script src="/static/js/home.js"></script>
</body>
</html>