export class AudioPlayer {
  constructor() {
    this.volume = 1.0; // Set default volume to 50%
    this.loop = false; // Set default loop to false
    this.muted = false; // Set default muted to false
    this.playbackRate = 1; // Set default playback rate to 1
    this.currentTime = 0; // Set default current time to 0
    // this.audio.src = ''; // Set default source to empty
    this.ids = {
      audio: "piped:audio",
      currentTime: "currentTime",
      totalTime: "totalTime",
      progress: "audio_progress",
    };
    this.uiElement = null;
    this.ui_CurrentTime = null;
    this.ui_TotalTime = null;
    this.ui_ProgressBar = null;
  }

  // init looks the elements up again. The <audio> element stays put, but the
  // controls only exist on the playlist page, so this runs after every page
  // load; on other pages the controls are simply null.
  init() {
    this.uiElement = document.getElementById(this.ids.audio);
    this.ui_CurrentTime = document.getElementById(this.ids.currentTime);
    this.ui_TotalTime = document.getElementById(this.ids.totalTime);
    this.ui_ProgressBar = document.getElementById(this.ids.progress);
    if (!this.uiElement) {
      console.error("Audio element not found.");
    }
  }

  updateTime() {
    if (this.ui_TotalTime && this.ui_CurrentTime && this.ui_ProgressBar) {
      const { minutes, seconds } = this.toMinutes(this.uiElement.duration);
      const { minutes: mCurrent, seconds: sCurrent } = this.toMinutes(this.uiElement.currentTime);
      if(isNaN(minutes) || isNaN(seconds)) {
        console.error("Invalid duration.");
        return;
      }
      this.ui_CurrentTime.innerHTML = `${mCurrent}:${parseFloat(sCurrent * 0.01).toFixed(2).toString().replace("0.", "")}`;
      this.ui_ProgressBar.value = (this.uiElement.currentTime / this.uiElement.duration) * 100;
      this.ui_TotalTime.innerHTML = `${minutes}:${parseFloat(seconds * 0.01).toFixed(2).toString().replace("0.", "")}`;
    }
  }

  toMinutes(totalSeconds) {
    const minutes = Math.floor(totalSeconds / 60);
    const remainingSeconds = totalSeconds % 60;
    return { minutes, seconds: remainingSeconds };
  }

  togglePlay() {
    if (window.pipedPlayer.state === "PAUSED") {
      window.pipedPlayer.resume();
    } else {
      window.pipedPlayer.pause();
    }
  }

  RR(){
    // Rewind
    this.uiElement.currentTime = this.uiElement.currentTime - 10;
  }
  FF(){
    // Fast Forward
    this.uiElement.currentTime = this.uiElement.currentTime + 10;
  }

  toggleMute() {
    if (this.uiElement.volume > 0) {
      this.uiElement.volume = 0;
    } else {
      this.uiElement.volume = this.volume;
    }
  }
}
