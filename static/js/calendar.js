import { fetchAPI } from './api.js';

let selectedDate = null;

export function generateContributionGraph() {
    const container = document.getElementById('contributionGraph');
    container.innerHTML = '';

    const today = new Date();
    const oneYearAgo = new Date(today.getFullYear() - 1, today.getMonth(), today.getDate());

    fetchAPI(`/api/contribution?start=${oneYearAgo.toISOString().split('T')[0]}&end=${today.toISOString().split('T')[0]}`)
        .then(data => {
            if (!data || !Array.isArray(data)) {
                container.innerHTML = '<p class="text-gray-400">Erro ao carregar dados</p>';
                return;
            }

            const dataMap = new Map();
            data.forEach(day => {
                dataMap.set(day.date, day);
            });

            for (let d = oneYearAgo; d <= today; d.setDate(d.getDate() + 1)) {
                const dateStr = d.toISOString().split('T')[0];
                const cell = document.createElement('div');
                cell.className = 'contribution-cell cursor-pointer';
                cell.dataset.date = dateStr;

                const dayData = dataMap.get(dateStr);
                let colorClass = 'color-0';

                if (dayData && dayData.avg_temp !== null) {
                    const temp = dayData.avg_temp;
                    if (temp < 15) colorClass = 'color-1';
                    else if (temp < 20) colorClass = 'color-2';
                    else if (temp < 25) colorClass = 'color-3';
                    else if (temp < 30) colorClass = 'color-4';
                    else colorClass = 'color-5';
                }

                cell.classList.add(colorClass);
                cell.onclick = () => openDayModal(dateStr);

                cell.addEventListener('mouseenter', (e) => showTooltip(e, dayData, dateStr));
                cell.addEventListener('mouseleave', hideTooltip);

                container.appendChild(cell);
            }
        });
}

export function openDayModal(date) {
    selectedDate = date;
    document.getElementById('dayModal').classList.remove('hidden');
    document.getElementById('dayModal').classList.add('flex');
    document.getElementById('modalDate').textContent = new Date(date).toLocaleDateString('pt-BR');
    document.getElementById('modalContent').innerHTML = '<p class="text-gray-400">Carregando...</p>';

    fetchAPI(`/api/day/${date}`)
        .then(data => {
            const contentDiv = document.getElementById('modalContent');
            
            if (!data || data.message) {
                contentDiv.innerHTML = '<p class="text-gray-400">Sem dados disponíveis para este dia</p>';
                return;
            }

            const formattedDate = new Date(data.date + 'T00:00:00').toLocaleDateString('pt-BR', {
                weekday: 'long',
                year: 'numeric',
                month: 'long',
                day: 'numeric'
            });

            contentDiv.innerHTML = `
                <div class="space-y-3">
                    <p class="text-gray-300 text-sm">${formattedDate}</p>
                    <div class="grid grid-cols-2 gap-4">
                        <div class="bg-gray-700 rounded p-3">
                            <p class="text-xs text-gray-400 mb-1">Média</p>
                            <p class="text-2xl font-bold text-green-400">${data.avg_temp.toFixed(1)}°C</p>
                        </div>
                        <div class="bg-gray-700 rounded p-3">
                            <p class="text-xs text-gray-400 mb-1">Máxima</p>
                            <p class="text-2xl font-bold text-red-400">${data.max_temp.toFixed(1)}°C</p>
                        </div>
                        <div class="bg-gray-700 rounded p-3">
                            <p class="text-xs text-gray-400 mb-1">Mínima</p>
                            <p class="text-2xl font-bold text-blue-400">${data.min_temp.toFixed(1)}°C</p>
                        </div>
                        <div class="bg-gray-700 rounded p-3">
                            <p class="text-xs text-gray-400 mb-1">Umidade</p>
                            <p class="text-2xl font-bold text-purple-400">${data.avg_humidity.toFixed(0)}%</p>
                        </div>
                    </div>
                </div>
            `;
        });
}

export function closeModal() {
    document.getElementById('dayModal').classList.add('hidden');
    document.getElementById('dayModal').classList.remove('flex');
    selectedDate = null;
}

export function editFeeling() {
    if (!selectedDate) return;
    document.getElementById('feelingDate').value = selectedDate;
    closeModal();
    document.getElementById('feelingForm').scrollIntoView({ behavior: 'smooth' });
}

function showTooltip(event, dayData, dateStr) {
    const tooltip = document.getElementById('tooltip');
    
    if (!dayData || dayData.avg_temp === null) {
        tooltip.innerHTML = `
            <div class="date">${new Date(dateStr + 'T00:00:00').toLocaleDateString('pt-BR')}</div>
            <div class="metric">Sem dados</div>
        `;
    } else {
        tooltip.innerHTML = `
            <div class="date">${new Date(dateStr + 'T00:00:00').toLocaleDateString('pt-BR')}</div>
            <div class="metric"><span>Temp Média:</span><span class="text-green-400">${dayData.avg_temp.toFixed(1)}°C</span></div>
            <div class="metric"><span>Máxima:</span><span class="text-red-400">${dayData.max_temp.toFixed(1)}°C</span></div>
            <div class="metric"><span>Mínima:</span><span class="text-blue-400">${dayData.min_temp.toFixed(1)}°C</span></div>
            <div class="metric"><span>Umidade:</span><span class="text-purple-400">${dayData.avg_humidity.toFixed(0)}%</span></div>
        `;
    }

    tooltip.classList.remove('hidden');

    const rect = event.target.getBoundingClientRect();
    const tooltipRect = tooltip.getBoundingClientRect();
    
    let left = rect.left + rect.width / 2 - tooltipRect.width / 2;
    let top = rect.top - tooltipRect.height - 8;

    if (left < 10) left = 10;
    if (left + tooltipRect.width > window.innerWidth - 10) {
        left = window.innerWidth - tooltipRect.width - 10;
    }
    if (top < 10) {
        top = rect.bottom + 8;
    }

    tooltip.style.left = `${left}px`;
    tooltip.style.top = `${top}px`;
}

function hideTooltip() {
    document.getElementById('tooltip').classList.add('hidden');
}
