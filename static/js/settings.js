import { fetchAPI } from './api.js';

export function openSettingsModal() {
    document.getElementById('settingsModal').classList.remove('hidden');
    document.getElementById('settingsModal').classList.add('flex');
    loadSettings();
}

export function closeSettingsModal() {
    document.getElementById('settingsModal').classList.add('hidden');
    document.getElementById('settingsModal').classList.remove('flex');
}

function loadSettings() {
    fetchAPI('/api/settings')
        .then(data => {
            if (data && data.sensorIntervalMinutes) {
                document.getElementById('sensorInterval').value = data.sensorIntervalMinutes;
            }
        });
}

export function initSettingsForm() {
    document.getElementById('settingsForm').addEventListener('submit', function(e) {
        e.preventDefault();
        
        const sensorInterval = document.getElementById('sensorInterval').value;
        
        fetchAPI('/api/settings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                sensorInterval: parseInt(sensorInterval)
            })
        })
        .then(data => {
            if (data && data.sensorIntervalMinutes) {
                return fetchAPI('/api/admin/restart-sensor', {
                    method: 'POST'
                });
            }
            return null;
        })
        .then(response => {
            if (response && response.message) {
                alert('Configurações salvas e sensor reiniciado com sucesso!');
                closeSettingsModal();
            } else {
                alert('Configurações salvas com sucesso!');
                closeSettingsModal();
            }
        });
    });
}
