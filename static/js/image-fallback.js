function showImageFallback(image) {
    if (image.naturalWidth > 0) {
        return;
    }

    image.hidden = true;

    const imageContainer = image.closest(
        ".card-image-link, .detail-image, .seatmap-section a"
    );

    if (imageContainer) {
        imageContainer.classList.add("image-failed");
    }
}

document.querySelectorAll("img").forEach((image) => {
    image.addEventListener("error", () => showImageFallback(image), { once: true });

    if (image.complete && image.naturalWidth === 0) {
        showImageFallback(image);
    }
});