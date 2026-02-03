import { loadHistory } from './history.js';
import { setRating } from './feeling.js';
import { closeModal, editFeeling } from './calendar.js';
import { openFeelingModal, closeFeelingModal } from './feeling.js';
import { openSettingsModal, closeSettingsModal } from './settings.js';

window.loadHistory = loadHistory;
window.setRating = setRating;
window.closeModal = closeModal;
window.editFeeling = editFeeling;
window.openFeelingModal = openFeelingModal;
window.closeFeelingModal = closeFeelingModal;
window.openSettingsModal = openSettingsModal;
window.closeSettingsModal = closeSettingsModal;
