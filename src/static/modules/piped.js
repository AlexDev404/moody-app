export class Piped {
  constructor() {
    this.apiUrl = "https://pipedapi.wireway.ch";
    // this.apiUrl = "https://pipedapi.ducks.party";
    // this.apiUrl = "https://pipedapi.kavin.rocks";
    this.playerId = "piped:audio";
    this.playerElement = null;
    this.playerInstance = null;
    this.volume = 1.0; // Default volume
    this.spinner = "piped:spinner";
    this.loadingSpinner = document.createElement("div");
    this.loadingSpinner.setAttribute("data-uk-spinner", "");
    // Player states
    this.STATE = {
      LOADING: "LOADING",
      ERROR: "ERROR",
      STOPPED: "STOPPED",
      PLAYING: "PLAYING",
      PAUSED: "PAUSED",
    };
    this.state = this.STATE.STOPPED; // Default state
  }

  init() {
    this.playerElement = document.getElementById(this.playerId);
    this.spinner = document.getElementById(this.spinner);
    if (!this.playerElement || !this.spinner) {
      console.error(
        `Player element with ID ${this.playerId} or spinner not found.`
      );
    }
    this.playerElement.volume = this.volume; // Set initial volume
    this.playerElement.addEventListener("ended", () => {
      this.updateState(this.STATE.STOPPED); // Update state to STOPPED when video ends
    });

    if (!window._playerStateChangeListener) {
      console.log("piped:init");
      // Create new player instance and remove spinner when ready
      document.addEventListener("piped:external_playerStateChange", (e) => {
        this.updateState(e.detail.state);
      });
      this.playerElement.onwaiting = () => {
        this.updateState(this.STATE.LOADING);
      };

      this.playerElement.onloadeddata = () => {
        document.dispatchEvent(
          new CustomEvent("piped:external_playerStateChange", {
            detail: { state: "PLAYING" },
          })
        );
      };

      this.playerElement.onplaying = () => {
        setTimeout(() => {
          this.updateState(this.STATE.PLAYING);
        }, 300);
      };

      window._playerStateChangeListener = true; // Set flag to prevent multiple listeners
    }
  }

  removeSpinner = () => {
    if (this.spinner.contains(this.loadingSpinner)) {
      this.spinner.removeChild(this.loadingSpinner);
    } else {
      console.warn("piped:spinner not found in the spinner element.");
    }
  };

  // Setter for volume
  setVolume(volume) {
    if (this.playerElement) {
      this.volume = volume;
      this.playerElement.volume = volume;
    } else {
      console.error("Player element not initialized.");
    }
  }

  updateState(newState) {
    if (!Object.values(this.STATE).includes(newState)) {
      console.error(`Invalid state: ${newState}`);
      return;
    }
    if (this.state === newState) {
      console.warn(`State is already ${newState}. No update needed.`);
      return;
    }
    console.log(`Player state updated from ${this.state} to ${newState}`);
    if (newState === this.STATE.LOADING) {
      this.spinner.appendChild(this.loadingSpinner);
    }
    if (newState === this.STATE.PLAYING) {
      this.removeSpinner();
    }
    if (newState === this.STATE.ERROR) {
      this.removeSpinner();
      console.error("Error occurred while playing the video.");
    }
    if (newState === this.STATE.STOPPED) {
      this.stop();
    }

    document.dispatchEvent(
      new CustomEvent("piped:playerStateChange", {
        detail: { state: newState },
      })
    );
    this.state = newState;
  }

  play(searchQuery) {
    if (!this.playerElement) {
      console.error("Player element not initialized.");
      return;
    }

    // Show loading state and spinner
    this.updateState(this.STATE.LOADING);

    // Clean up any previous player instance
    if (this.playerInstance) {
      this.playerInstance.stop();
    }

    this.playerInstance = new PipedPlayer(this.playerElement, {
      autoplay: true,
      search: searchQuery,
      apiUrl: this.apiUrl,
    });

    this.playerInstance.play();
    this.updateState(this.STATE.PLAYING); // Update state to PLAYING
  }
  pause() {
    if (this.playerInstance) {
      this.playerInstance.pause();
      this.updateState(this.STATE.PAUSED); // Update state to PAUSED
    } else {
      console.error("Player instance not initialized.");
    }
  }
  stop() {
    if (this.playerInstance) {
      this.playerInstance.stop();
    } else {
      console.error("Player instance not initialized.");
    }
  }
  resume() {
    if (this.playerInstance) {
      this.playerInstance.resume();
      this.updateState(this.STATE.PLAYING); // Update state to PLAYING
    } else {
      console.error("Player instance not initialized.");
    }
  }
}

export class PipedPlayer {
  constructor(element, options) {
    this.element = element;
    this.options = options;
  }

  async play() {
    const searchQuery = this.options.search;
    const apiUrl = this.options.apiUrl;
    const url = `${apiUrl}/search?q=${encodeURIComponent(
      searchQuery
    )}&filter=all`;
    await fetch(url)
      .then((response) => response.json())
      .then(async (data) => {
        if (data?.items?.length > 0) {
          // Find the first stream-type item with duration > 0
          const video =
            data.items.find(
              (item) => item.type === "stream" && item.duration > 0
            ) || data.items[0]; // Fallback to first item if no suitable stream found
          this.element.src = video.url;
          const videoId = video.url.replace(/.*\/watch\?v=/, "");
          const res = await fetch(`${apiUrl}/streams/${videoId}`);
          const streamData = await res.json();
          // Find the first audio stream with mimetype "audio/webm"
          const webmAudio = streamData?.audioStreams?.find((stream) =>
            stream.mimeType?.includes("audio/webm")
          );
          const streamUrl =
            webmAudio?.url || streamData?.audioStreams?.[0]?.url;
          if (streamUrl) {
            this.element.src = streamUrl;
            this.element.play();
          } else {
            console.error("No suitable audio stream found.");
          }
        } else {
          console.error("No videos found for the search query.");
        }
      })
      .catch((error) => {
        console.error("Error fetching video:", error);
      });
    if (this.options.autoplay) {
      this.element.play();
    }
  }

  pause() {
    if (this.element) {
      this.element.pause();
    } else {
      console.error("Player element not initialized.");
    }
  }

  stop() {
    if (this.element) {
      this.element.pause();
      this.element.currentTime = 0;
      // Destroy the player instance
      this.element.src = ""; // Clear the source
    } else {
      console.error("Player element not initialized.");
    }
  }

  resume() {
    if (this.element) {
      this.element.play();
    } else {
      console.error("Player element not initialized.");
    }
  }
}
