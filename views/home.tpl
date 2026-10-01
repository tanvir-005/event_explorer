<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#fafbf8">
    <title>Discover your next plan | Event Explorer</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body class="page-home">
    {{template "layouts/header.tpl" .}}

    <main id="main" class="container home-page">
        <section class="home-hero" aria-labelledby="home-title">
            <div class="hero-copy">
                <p class="eyebrow"><span class="eyebrow-mark">01</span> YOUR CITY, IN THE MOMENT</p>
                <h1 id="home-title">Find your next<br><em>out.</em></h1>
                <p class="hero-subtitle">Good music. Big games. A reason to make plans.<br>See what's happening in your city.</p>
                <div class="hero-tags" aria-label="Event categories">
                    <span>Music</span>
                    <span>Sports</span>
                    <span>Live near you</span>
                </div>
            </div>
            <figure class="hero-photo">
                <img src="/static/img/concert.jpg" alt="A live music crowd under bright stage lights">
                <figcaption><span>OUT THERE / TOGETHER</span><span>LIVE / LOCAL</span></figcaption>
            </figure>
        </section>

        <section class="search-panel" aria-labelledby="search-title">
            <div class="search-intro">
                <div>
                    <p class="eyebrow">START WITH A CITY</p>
                    <h2 id="search-title">Where are you going?</h2>
                </div>
                <p>Find live music and sports<br>where you want to be.</p>
            </div>

            <form id="city-form" class="search-form">
                <div class="input-shell">
                    <label for="city-input">City</label>
                    <input
                        type="text"
                        id="city-input"
                        placeholder="Try New York or London"
                        autocomplete="off"
                        role="combobox"
                        aria-autocomplete="list"
                        aria-controls="suggestions"
                        aria-expanded="false"
                    >
                    <div id="suggestions" class="suggestions" role="listbox" aria-label="City suggestions"></div>
                </div>
                <input type="hidden" id="place-id">
                <input type="hidden" id="session-token">
                <button class="button button-primary search-button" type="submit">Explore events <span aria-hidden="true">&rarr;</span></button>
            </form>
            <div class="search-meta">
                <p id="search-status" role="status" aria-live="polite"></p>
                <p id="error" role="alert"></p>
                <p class="google-attribution">Location suggestions by Google</p>
            </div>
        </section>

        <section class="home-bottom" aria-label="Explore event categories">
            <p class="mini-label">MAKE A NIGHT OF IT</p>
            <div class="category-line">
                <h2>Pick your kind of live.</h2>
                <p>Two good places to start <span aria-hidden="true">&rarr;</span></p>
            </div>
            <div class="category-list">
                <div><span class="category-index">01</span><span class="category-name">Music</span><span class="category-note">From intimate rooms to arena nights</span></div>
                <div><span class="category-index">02</span><span class="category-name">Sports</span><span class="category-note">The home crowd, the big moment</span></div>
            </div>
        </section>
    </main>

    {{template "layouts/footer.tpl" .}}

    <script src="/static/js/home.js"></script>
</body>
</html>