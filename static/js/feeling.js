import { generateContributionGraph } from './calendar.js';

let currentRating = 0;

export function setRating(rating) {
    currentRating = rating;
    document.getElementById('feelingRating').value = rating;
    const stars = document.querySelectorAll('#ratingStars span');
    stars.forEach((star, index) => {
        if (index < rating) {
            star.classList.remove('text-gray-500');
            star.classList.add('text-yellow-400');
        } else {
            star.classList.remove('text-yellow-400');
            star.classList.add('text-gray-500');
        }
    });
}

export function initFeelingForm() {
    document.getElementById('feelingDate').valueAsDate = new Date();
    setRating(0);

    document.getElementById('feelingForm').addEventListener('htmx:afterRequest', function(evt) {
        if (evt.detail.xhr.status === 200) {
            alert('Sensação salva com sucesso!');
            document.getElementById('feelingForm').reset();
            setRating(0);
            generateContributionGraph();
        }
    });
}
