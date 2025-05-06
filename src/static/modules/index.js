import { AudioPlayer } from "./audio.js";
import { addCSRFToForm, fetchWithCSRF, getCSRFToken } from "./csrf.js";
import { Piped, PipedPlayer } from "./piped.js";

window.Piped = Piped;
window.PipedPlayer = PipedPlayer;
window.AudioPlayer = AudioPlayer;

// Export CSRF utilities to global scope for use in inline scripts
window.fetchWithCSRF = fetchWithCSRF;
window.getCSRFToken = getCSRFToken;
window.addCSRFToForm = addCSRFToForm;
