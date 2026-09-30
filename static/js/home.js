const input = document.getElementById("city-input");
const suggestions = document.getElementById("suggestions");
const placeIdInput = document.getElementById("place-id");
const sessionTokenInput = document.getElementById("session-token");
const form = document.getElementById("city-form");
const error = document.getElementById("error");

let debounceTimer;

function newSessionToken() {
    return crypto.randomUUID();
}

input.addEventListener("input", () => {
    clearTimeout(debounceTimer);

    placeIdInput.value = "";
    error.textContent = "";

    sessionTokenInput.value = newSessionToken();

    const value = input.value.trim();

    if (!value) {
        suggestions.innerHTML = "";
        return;
    }

    debounceTimer = setTimeout(async () => {
        try {
            const token = sessionTokenInput.value;

            const response = await fetch(
                `/api/locations/autocomplete?input=${encodeURIComponent(value)}&sessionToken=${encodeURIComponent(token)}`
            );

            const data = await response.json();

            suggestions.innerHTML = "";

            if (!response.ok) {
                error.textContent = data.error || "Unable to fetch suggestions.";
                return;
            }

            for (const suggestion of data.suggestions) {
                const button = document.createElement("button");

                button.type = "button";
                button.textContent = suggestion.text;

                button.addEventListener("click", () => {
                    input.value = suggestion.text;
                    placeIdInput.value = suggestion.placeId;
                    suggestions.innerHTML = "";
                });

                suggestions.appendChild(button);
            }
        } catch {
            error.textContent = "Unable to fetch suggestions.";
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

    try {
        const response = await fetch(
            `/api/locations/${encodeURIComponent(placeId)}?sessionToken=${encodeURIComponent(sessionToken)}`
        );

        const data = await response.json();

        if (!response.ok) {
            error.textContent = data.error || "Unable to resolve city.";
            return;
        }

        window.location.href =
            `/events?city=${encodeURIComponent(data.city)}&countryCode=${encodeURIComponent(data.countryCode)}`;
    } catch {
        error.textContent = "Unable to resolve city.";
    }
});