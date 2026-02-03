import { fetchAPI } from './api.js';

let historyChart = null;

export function loadHistory() {
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;

    if (!startDate || !endDate) {
        alert('Selecione as datas');
        return;
    }

    fetchAPI(`/api/history?start=${new Date(startDate).toISOString()}&end=${new Date(endDate).toISOString()}`)
        .then(data => {
            if (!data || !Array.isArray(data)) {
                return;
            }

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
        });
}
