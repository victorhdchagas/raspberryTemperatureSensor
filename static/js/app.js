let historyChart = null;
let selectedDate = null;

async function fetchAPI(endpoint, options = {}) {
    try {
        const response = await fetch(endpoint, options);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
        return await response.json();
    } catch (error) {
        console.error(`Error fetching ${endpoint}:`, error);
        return null;
    }
}

async function loadHistory() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;

    if (!startDate || !endDate) {
        alert('Selecione as datas');
        return;
    }

    const data = await fetchAPI(`/api/history?start=${new Date(startDate).toISOString()}&end=${new Date(endDate).toISOString()}`);
    if (data && Array.isArray(data)) {
        const labels = data.map(m => new Date(m.timestamp).toLocaleString('pt-BR'));
        const temps = data.map(m => m.temp);
        const humidity = data.map(m => m.humidity);

        if (historyChart) {
            historyChart.destroy();
        }

        const ctx = document.getElementById('historyChart').getContext('2d');
        historyChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Temperatura (°C)',
                        data: temps,
                        borderColor: 'rgb(74, 222, 128)',
                        backgroundColor: 'rgba(74, 222, 128, 0.1)',
                        tension: 0.1
                    },
                    {
                        label: 'Umidade (%)',
                        data: humidity,
                        borderColor: 'rgb(96, 165, 250)',
                        backgroundColor: 'rgba(96, 165, 250, 0.1)',
                        tension: 0.1,
                        yAxisID: 'y1'
                    }
                ]
            },
            options: {
                responsive: true,
                interaction: {
                    mode: 'index',
                    intersect: false,
                },
                scales: {
                    x: {
                        ticks: {
                            color: '#9ca3af',
                            maxTicksLimit: 10
                        },
                        grid: {
                            color: '#374151'
                        }
                    },
                    y: {
                        type: 'linear',
                        display: true,
                        position: 'left',
                        ticks: {
                            color: '#9ca3af'
                        },
                        grid: {
                            color: '#374151'
                        }
                    },
                    y1: {
                        type: 'linear',
                        display: true,
                        position: 'right',
                        grid: {
                            drawOnChartArea: false,
                        },
                        ticks: {
                            color: '#9ca3af'
                        }
                    },
                },
                plugins: {
                    legend: {
                        labels: {
                            color: '#e5e7eb'
                        }
                    }
                }
            }
        });
    }
}

function generateContributionGraph() {
    const container = document.getElementById('contributionGraph');
    container.innerHTML = '';

    const today = new Date();
    const oneYearAgo = new Date(today.getFullYear() - 1, today.getMonth(), today.getDate());

    for (let d = oneYearAgo; d <= today; d.setDate(d.getDate() + 1)) {
        const cell = document.createElement('div');
        cell.className = 'contribution-cell cursor-pointer hover:opacity-80';
        cell.dataset.date = d.toISOString().split('T')[0];

        const random = Math.random();
        let colorClass = 'color-0';
        if (random < 0.2) colorClass = 'color-1';
        else if (random < 0.4) colorClass = 'color-2';
        else if (random < 0.6) colorClass = 'color-3';
        else if (random < 0.8) colorClass = 'color-4';

        cell.classList.add(colorClass);
        cell.onclick = () => openDayModal(d.toISOString().split('T')[0]);
        container.appendChild(cell);
    }
}

function openDayModal(date) {
    selectedDate = date;
    document.getElementById('dayModal').classList.remove('hidden');
    document.getElementById('dayModal').classList.add('flex');
    document.getElementById('modalDate').textContent = new Date(date).toLocaleDateString('pt-BR');
    document.getElementById('modalContent').innerHTML = '<p class="text-gray-400">Carregando...</p>';
}

function closeModal() {
    document.getElementById('dayModal').classList.add('hidden');
    document.getElementById('dayModal').classList.remove('flex');
    selectedDate = null;
}

function editFeeling() {
    if (!selectedDate) return;
    document.getElementById('feelingDate').value = selectedDate;
    closeModal();
    document.getElementById('feelingForm').scrollIntoView({ behavior: 'smooth' });
}

function setRating(rating) {
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

document.getElementById('feelingForm').addEventListener('htmx:afterRequest', function(evt) {
    if (evt.detail.xhr.status === 200) {
        alert('Sensação salva com sucesso!');
        document.getElementById('feelingForm').reset();
        setRating(0);
    }
});

document.getElementById('feelingDate').valueAsDate = new Date();
const endDate = new Date();
const startDate = new Date(endDate.getTime() - 24 * 60 * 60 * 1000);
document.getElementById('startDate').value = startDate.toISOString().slice(0, 16);
document.getElementById('endDate').value = endDate.toISOString().slice(0, 16);

generateContributionGraph();
loadHistory();
