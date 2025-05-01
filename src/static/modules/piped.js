export class Piped {
  constructor() {
    this.apiUrl = "https://pipedapi.ducks.party";
    this.playerId = "piped-audio";
    this.playerElement = null;
    this.playerInstance = null;
    this.volume = 1.0; // Default volume
  }

  init() {
    this.playerElement = document.getElementById(this.playerId);
    if (!this.playerElement) {
      console.error(`Player element with ID ${this.playerId} not found.`);
    }
    this.playerElement.volume = this.volume; // Set initial volume
    this.playerElement.addEventListener("ended", () => {
      this.audio.currentTime = 0; // Reset current time to 0 when audio ends
    });
  }

  // Setter for volume
  setVolume(volume) {
    if (this.playerElement) {
      this.volume = volume;
      this.playerElement.volume = volume;
    } else {
      console.error("Player element not initialized.");
    }
  }

  play(searchQuery) {
    if (!this.playerElement) {
      console.error("Player element not initialized.");
      return;
    }

    this.playerInstance = new PipedPlayer(this.playerElement, {
      autoplay: true,
      search: searchQuery,
      apiUrl: this.apiUrl,
    });

    this.playerInstance.play();
  }
  pause() {
    if (this.playerInstance) {
      this.playerInstance.pause();
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
