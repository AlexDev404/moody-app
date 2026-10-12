import { AudioPlayer } from "./audio.js";
import { Piped, PipedPlayer } from "./piped.js";

window.Piped = Piped;
window.PipedPlayer = PipedPlayer;

// One player for the whole visit: the <audio> element sits outside #content,
// so music keeps playing while folio swaps pages. The controls live on the
// playlist page, so they're looked up again after every page load.
window.pipedPlayer = new Piped();
window.AudioPlayer = new AudioPlayer();

function bindPlayer() {
  window.pipedPlayer.init();
  window.AudioPlayer.init();
}
document.addEventListener("appLoad", bindPlayer);
if (window.folio?.loaded) bindPlayer();
