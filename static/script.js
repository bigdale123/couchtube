const IFRAME_API_URL = 'https://www.youtube.com/iframe_api';
const BUFFERING_TIMEOUT = 3000;
const CHANNELS_ENDPOINT = '/channels';
const CURRENT_VIDEO_ENDPOINT = '/current-video';
const VOLUME_STEPS = 5;
const VOLUME_BAR_TIMEOUT = 2000;

const ICONS = {
  power: '/assets/icons/power.svg',
  volume_muted: '/assets/icons/volume_muted.svg',
  volume_high: '/assets/icons/volume_high.svg'
};

class YouTubePlayer {
  constructor(playerElementId) {
    this.player = null;
    this.playerReady = false;
    this.playerElementId = playerElementId;
    this.channelsUrl = CHANNELS_ENDPOINT;
    this.currentVideoEndpoint = CURRENT_VIDEO_ENDPOINT;
    this.channels = [];
    this.currentChannel = null;
    this.videoTitle = '';
    this.isPlaying = false;
    this.isMuted = true;
    this.hasInteracted = false;
    this.volumeBarTimeoutId = null;
    this.currentVideo = null;

    this.loadYouTubeAPI();
    this.loadChannels().then(() => {
      this.addControlListeners();
    });
  }

  loadYouTubeAPI() {
    const tag = document.createElement('script');
    tag.src = IFRAME_API_URL;
    const firstScriptTag = document.getElementsByTagName('script')[0];
    firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

    window.onYouTubeIframeAPIReady = () => this.onYouTubeIframeAPIReady();
  }

  onYouTubeIframeAPIReady() {
    this.player = new YT.Player(this.playerElementId, {
      width: '100%',
      height: '100%',
      autoplay: 1,
      events: {
        onReady: (event) => this.onPlayerReady(event),
        onStateChange: (event) => this.onPlayerStateChange(event)
      },
      playerVars: {
        mute: 1,
        controls: 0,
        modestbranding: 1,
        disablekb: 1,
        fs: 0,
        iv_load_policy: 3,
        rel: 0,
        enablejsapi: 1,
        autoplay: 1
      }
    });
  }

  onPlayerReady(event) {
    this.playerReady = true;
    if (this.channels.length > 0) {
      this.loadChannelVideo(this.channels[0]);
    } else {
      console.warn('No channels available to load.');
    }
  }

  onPlayerStateChange({ data, target }) {
    this.videoTitle = target.videoTitle || this.videoTitle;

    if (data === YT.PlayerState.UNSTARTED) {
      this.playVideo();
    }

    if (data === YT.PlayerState.PLAYING) {
      const intervalId = setInterval(async () => {
        const currentTime = this.player.getCurrentTime();

        // If current time reaches or exceeds the end time, load the next video
        if (currentTime >= this.currentVideo.segmentEnd) {
          clearInterval(intervalId);
          await this.loadChannelVideo(this.currentChannel, this.currentVideo);
        }
      }, 1000);
    }
  }

  async loadChannels() {
    try {
      const res = await fetch(this.channelsUrl);
      const data = await res.json();
      this.channels = data.channels || [];
    } catch (error) {
      console.error('Failed to load channels:', error);
    }
  }

  async getCurrentVideo(channelId, videoId) {
    try {
      // non-null videoId means we're looking for the next video
      const url = `${this.currentVideoEndpoint}?channel-id=${channelId}${
        videoId ? `&video-id=${videoId}` : ''
      }`;
      const res = await fetch(url);
      const data = await res.json();
      if (data.video) return data.video;
      return null;
    } catch (error) {
      console.error('Failed to load channels:', error);
    }
  }

  playVideo() {
    if (
      this.playerReady &&
      this.player.getPlayerState() !== YT.PlayerState.PLAYING
    ) {
      if (!this.hasInteracted) this.player.mute();
      this.showBuffering();

      setTimeout(() => {
        if (this.hasInteracted && !this.isMuted) this.player.unMute();
        this.hideBuffering();
        this.player.playVideo();
      }, BUFFERING_TIMEOUT);
    }
  }

  pauseVideo() {
    if (
      this.playerReady &&
      this.player.getPlayerState() === YT.PlayerState.PLAYING
    ) {
      this.player.pauseVideo();
      this.showBuffering();
    }
  }

  updateIcon(iconId, iconSrc, isActive) {
    const iconElement = document.querySelector(`#${iconId} .control-icon`);
    if (iconElement) {
      iconElement.src = iconSrc;
      iconElement.classList.toggle('red', isActive);
    }
  }

  toggleMute() {
    if (this.playerReady) {
      this.isMuted = !this.isMuted;
      this.player[this.isMuted ? 'mute' : 'unMute']();
      this.updateIcon(
        'control-mute',
        this.isMuted ? ICONS.volume_muted : ICONS.volume_high,
        this.isMuted
      );
    }
  }

  async loadChannelVideo(channel, currentVideo = null) {
    const videoToBePlayed = await this.getCurrentVideo(
      channel.id,
      currentVideo?.id
    );

    if (!videoToBePlayed) {
      console.error('No video found for channel:', channel);
      return;
    }

    this.currentVideo = videoToBePlayed;
    const videoUrl = videoToBePlayed.url;
    const videoId = this.extractVideoId(videoUrl);

    if (videoId && this.playerReady) {
      this.pauseVideo();
      this.player.cueVideoById({
        videoId,
        startSeconds: videoToBePlayed.segmentStart
      });

      this.playVideo();
      this.currentChannel = channel;
    } else {
      console.error('Player not ready or invalid video ID:', videoId);
    }
  }

  extractVideoId(url) {
    const match = url.match(/(?:v=|\/)([0-9A-Za-z_-]{11})/);
    return match ? match[1] : null;
  }

  changeChannelByOffset(offset) {
    const currentIndex = this.channels.findIndex(
      (channel) => channel.id === this.currentChannel.id
    );
    const newIndex =
      (currentIndex + offset + this.channels.length) % this.channels.length;
    this.loadChannelVideo(this.channels[newIndex]);
  }

  nextChannel() {
    this.changeChannelByOffset(1);
  }

  previousChannel() {
    this.changeChannelByOffset(-1);
  }

  showBuffering() {
    document.querySelector('#buffer-gif')?.classList.add('active');
  }

  hideBuffering() {
    document.querySelector('#buffer-gif')?.classList.remove('active');
  }

  updateVolumeBar(currentVolume) {
    const volumeBar = document.querySelector('#volume-bar');
    const maxBars = 100 / VOLUME_STEPS;
    const currentStep = Math.ceil(currentVolume / VOLUME_STEPS);

    volumeBar.classList.add('active');
    volumeBar.innerHTML = Array.from(
      { length: maxBars },
      (_, index) =>
        `<div class="volume-bar-step ${
          index < currentStep ? 'active' : ''
        }"></div>`
    ).join('');

    clearTimeout(this.volumeBarTimeoutId);
    this.volumeBarTimeoutId = setTimeout(
      () => volumeBar.classList.remove('active'),
      VOLUME_BAR_TIMEOUT
    );
  }

  adjustVolume(increase) {
    if (this.playerReady && increase !== undefined) {
      const currentVolume = this.player.getVolume();
      const adjustment = increase ? VOLUME_STEPS : -VOLUME_STEPS;
      const newVolume = Math.min(Math.max(currentVolume + adjustment, 0), 100); // cap between 0 and 100

      this.player.setVolume(newVolume);
      this.updateVolumeBar(newVolume);
      this.player.unMute();
    }
  }

  volumeUp() {
    this.adjustVolume(true);
  }

  volumeDown() {
    this.adjustVolume(false);
  }

  turnOff() {
    this.pauseVideo();
    this.updateIcon('control-power', ICONS.power, true);
  }

  turnOn() {
    this.playVideo();
    this.updateIcon('control-power', ICONS.power, false);
  }

  addControlListeners() {
    const controls = {
      power: () => {
        if (this.player.getPlayerState() === YT.PlayerState.PLAYING) {
          this.turnOff();
        } else {
          this.turnOn();
        }
      },
      chup: () => this.nextChannel(),
      chdown: () => this.previousChannel(),
      mute: () => this.toggleMute(),
      volup: () => this.volumeUp(),
      voldown: () => this.volumeDown()
    };

    for (const [control, handler] of Object.entries(controls)) {
      document
        .querySelector(`#control-${control}`)
        ?.addEventListener('click', handler);
    }

    document.querySelector('#controls')?.addEventListener('click', () => {
      this.hasInteracted = true;
    });

    document
      .querySelector('#control-fullscreen')
      ?.addEventListener('click', () => {
        const playerElement = document.querySelector('#player');

        console.log(playerElement);

        const requestFullScreen =
          playerElement.requestFullscreen ||
          playerElement.mozRequestFullScreen ||
          playerElement.webkitRequestFullScreen ||
          playerElement.msRequestFullscreen;

        if (requestFullScreen) {
          requestFullScreen.call(playerElement);
        }
      });
  }
}

const youtubePlayer = new YouTubePlayer('player');
