const input = document.getElementById("city-input");
const suggestions = document.getElementById("suggestions");
const placeIdInput = document.getElementById("place-id");
const sessionTokenInput = document.getElementById("session-token");
const form = document.getElementById("city-form");
const error = document.getElementById("error");

let debounceTimer;
let requestSequence = 0;

const status = document.getElementById("search-status");
const searchButton = form.querySelector("button[type='submit']");

function newSessionToken() {
    return crypto.randomUUID();
}

input.addEventListener("input", () => {
    clearTimeout(debounceTimer);
    requestSequence += 1;

    placeIdInput.value = "";
    error.textContent = "";
    status.textContent = "";
    suggestions.innerHTML = "";

    sessionTokenInput.value = newSessionToken();

    const value = input.value.trim();
    input.setAttribute("aria-expanded", String(Boolean(value)));

    if (!value) {
        return;
    }

    debounceTimer = setTimeout(async () => {
        const sequence = requestSequence;
        status.textContent = "Searching cities...";

        try {
            const token = sessionTokenInput.value;

            const response = await fetch(
                `/api/locations/autocomplete?input=${encodeURIComponent(value)}&sessionToken=${encodeURIComponent(token)}`
            );

            const data = await response.json();

            if (sequence !== requestSequence) {
                return;
            }

            suggestions.innerHTML = "";
            input.setAttribute("aria-expanded", "false");

            if (!response.ok) {
                error.textContent = data.error || "Unable to fetch suggestions.";
                status.textContent = "";
                return;
            }

            if (!data.suggestions.length) {
                status.textContent = "No matching cities found.";
                return;
            }

            for (const suggestion of data.suggestions) {
                const button = document.createElement("button");

                button.type = "button";
                button.setAttribute("role", "option");
                button.textContent = suggestion.text;

                button.addEventListener("click", () => {
                    input.value = suggestion.text;
                    placeIdInput.value = suggestion.placeId;
                    suggestions.innerHTML = "";
                    input.setAttribute("aria-expanded", "false");
                    status.textContent = "City selected. Explore events when you're ready.";
                });

                suggestions.appendChild(button);
            }
            input.setAttribute("aria-expanded", "true");
            status.textContent = `${data.suggestions.length} city suggestions available.`;
        } catch {
            if (sequence === requestSequence) {
                error.textContent = "Unable to fetch suggestions.";
                status.textContent = "";
                input.setAttribute("aria-expanded", "false");
            }
        }
    }, 300);
});

form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const placeId = placeIdInput.value;
    const sessionToken = sessionTokenInput.value;

    if (!placeId) {
        error.textContent = "Please select a city from the suggestions.";
        return;
    }

    searchButton.disabled = true;
    searchButton.setAttribute("aria-busy", "true");
    searchButton.firstChild.textContent = "Finding events ";
    error.textContent = "";

    try {
        const response = await fetch(
            `/api/locations/${encodeURIComponent(placeId)}?sessionToken=${encodeURIComponent(sessionToken)}`
        );

        const data = await response.json();

        if (!response.ok) {
            error.textContent = data.error || "Unable to resolve city.";
            status.textContent = "";
            return;
        }

        window.location.href =
            `/events?city=${encodeURIComponent(data.city)}&countryCode=${encodeURIComponent(data.countryCode)}`;
    } catch {
        error.textContent = "Unable to resolve city.";
    } finally {
        searchButton.disabled = false;
        searchButton.removeAttribute("aria-busy");
        searchButton.firstChild.textContent = "Explore events ";
    }
});