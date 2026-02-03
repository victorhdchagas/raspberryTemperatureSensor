import './global.js';
import { generateContributionGraph } from './calendar.js';
import { loadHistory } from './history.js';
import { initFeelingForm } from './feeling.js';

document.addEventListener('DOMContentLoaded', () => {
    const endDate = new Date();
    const startDate = new Date(endDate.getTime() - 24 * 60 * 60 * 1000);
    document.getElementById('startDate').value = startDate.toISOString().slice(0, 16);
    document.getElementById('endDate').value = endDate.toISOString().slice(0, 16);

    generateContributionGraph();
    loadHistory();
    initFeelingForm();
});
