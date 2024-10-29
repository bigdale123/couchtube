const IFRAME_API_URL = 'https://www.youtube.com/iframe_api';
const BUFFERING_TIMEOUT = 2000;
const CHANNELS_ENDPOINT = '/channels';

class YouTubePlayer {
  constructor(playerElementId) {
    this.player = null;
    this.playerReady = false;
    this.playerElementId = playerElementId;
    this.channelsUrl = CHANNELS_ENDPOINT;
    this.channels = [];
    this.currentChannel = null;
    this.videoTitle = '';
    this.isPlaying = false;
    this.isBuffering = true;

    this.loadYouTubeAPI();
    this.loadChannels();
    this.addControlListeners();
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
        controls: 0,
        modestbranding: 1,
        disablekb: 1,
        fs: 0,
        iv_load_policy: 3,
        rel: 0,
        enablejsapi: 1,
        loop: 0,
        autoplay: 1
      }
    });
  }

  onPlayerReady(event) {
    this.playerReady = true;
    if (this.channels.length > 0) {
      const initialChannel = this.channels[0];
      this.loadChannelVideo(initialChannel);
      this.currentChannel = initialChannel;
      event.target.playVideo();
    } else {
      console.warn('No channels available to load.');
    }
  }

  onPlayerStateChange(event) {
    const stateNames = {
      '-1': 'UNSTARTED',
      0: 'ENDED',
      1: 'PLAYING',
      2: 'PAUSED',
      3: 'BUFFERING',
      5: 'CUED'
    };

    console.log(event);

    console.log(`Current state: ${stateNames[event.data]}`);

    const videoTitle = event.target.videoTitle;
    if (videoTitle) this.videoTitle = videoTitle;

    if (event.data === YT.PlayerState.PLAYING) {
      this.isPlaying = true;

      setTimeout(() => this.deactivateBuffering(), BUFFERING_TIMEOUT);
      this.player.unMute();
    } else if (event.data === YT.PlayerState.BUFFERING) {
      this.activateBuffering();
    } else if (event.data === YT.PlayerState.PAUSED) {
      this.isPlaying = false;
      this.activateBuffering();
    } else if (event.data === YT.PlayerState.UNSTARTED) {
      this.playVideo();
    }
  }

  playVideo() {
    if (
      this.playerReady &&
      this.player.getPlayerState() !== YT.PlayerState.PLAYING
    ) {
      console.log('Playing video:', this.videoTitle);
      this.player.mute();
      this.player.playVideo();
    }
  }

  pauseVideo() {
    if (
      this.playerReady &&
      this.player.getPlayerState() === YT.PlayerState.PLAYING
    ) {
      this.player.pauseVideo();
    }
  }

  changeChannel(channel) {
    this.loadChannelVideo(channel);
    this.currentChannel = channel;
  }

  changeChannelByOffset(offset) {
    const currentIndex = this.channels.findIndex(
      (channel) => channel.id === this.currentChannel.id
    );
    const newIndex =
      (currentIndex + offset + this.channels.length) % this.channels.length;
    this.changeChannel(this.channels[newIndex]);
  }

  nextChannel() {
    this.changeChannelByOffset(1);
  }

  previousChannel() {
    this.changeChannelByOffset(-1);
  }

  addBufferingClass = () => {
    document.querySelector('#buffer-gif').classList.add('active');
  };

  removeBufferingClass = () => {
    document.querySelector('#buffer-gif').classList.remove('active');
  };

  activateBuffering = () => {
    if (!this.isBuffering) {
      this.isBuffering = true;
      this.addBufferingClass();
      this.player.mute();
    }
  };

  deactivateBuffering = () => {
    this.isBuffering = false;
    this.removeBufferingClass();
    this.player.unMute();
  };

  loadChannels() {
    fetch(this.channelsUrl)
      .then((res) => res.json())
      .then((res) => {
        this.channels = res.channels;
      });
  }

  loadChannelVideo(channel) {
    const videoUrl = channel.videos[0].url;
    const videoIdMatch = videoUrl.match(/(?:v=|\/)([0-9A-Za-z_-]{11})/);
    const videoId = videoIdMatch ? videoIdMatch[1] : null;
    console.log('Loading video ID:', videoId);

    if (videoId && this.playerReady) {
      this.player.loadVideoById(videoId);
    } else {
      console.error('Player not ready or invalid video ID:', videoId);
    }
  }

  addControlListeners() {
    document.querySelector('#control-power').addEventListener('click', () => {
      if (this.playerReady) {
        if (this.isPlaying) {
          this.pauseVideo();
        } else {
          this.playVideo();
        }
      }
    });

    document.querySelector('#control-chup').addEventListener('click', () => {
      this.nextChannel();
    });

    document.querySelector('#control-chdown').addEventListener('click', () => {
      this.previousChannel();
    });

    document.querySelector('#control-mute').addEventListener('click', () => {
      if (this.playerReady) {
        this.player.isMuted() ? this.player.unMute() : this.player.mute();
      }
    });
  }
}

const youtubePlayer = new YouTubePlayer('player');
