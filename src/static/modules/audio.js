export class AudioPlayer {
  constructor() {
    this.volume = 1.0; // Set default volume to 50%
    this.loop = false; // Set default loop to false
    this.muted = false; // Set default muted to false
    this.playbackRate = 1; // Set default playback rate to 1
    this.currentTime = 0; // Set default current time to 0
    // this.audio.src = ''; // Set default source to empty
    this.uiElement = "piped-audio"; // Set default UI element to null
    this.ui_CurrentTime = "currentTime"; // Set default UI current time to null
    this.ui_TotalTime = "totalTime"; // Set default UI current time to null
  }
  init() {
    this.ui_CurrentTime = document.getElementById(this.ui_CurrentTime);
    this.ui_TotalTime = document.getElementById(this.ui_TotalTime);
    if (!this.ui_CurrentTime || !this.ui_TotalTime) {
      console.error("UI elements not found.");
      return;
    }
    this.uiElement = document.getElementById(this.uiElement);
    if (!this.uiElement) {
      console.error("UI element not found.");
    }
  }

  updateTime() {
    if (this.ui_TotalTime) {
      const { minutes, seconds } = this.toMinutes(this.uiElement.duration);
      if(isNaN(minutes) || isNaN(seconds)) {
        console.error("Invalid duration.");
        return;
      }
      this.ui_TotalTime.innerHTML = `${minutes}:${Math.floor(seconds)}`;
    } else {
      console.error("UI element not initialized.");
    }
  }

  toMinutes(totalSeconds) {
    const minutes = Math.floor(totalSeconds / 60);
    const remainingSeconds = totalSeconds % 60;
    return { minutes, seconds: remainingSeconds };
  }
}
