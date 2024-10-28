const IFRAME_API_URL = 'https://www.youtube.com/iframe_api';
const BUFFERING_TIMEOUT = 2000;

class YouTubePlayer {
  channelsUrl = '/channels';

  constructor(playerElementId) {
    this.player = null;
    this.playerReady = false;
    this.playerElementId = playerElementId;

    this.loadYouTubeAPI();
    this.loadChannels();
    this.addControlListeners();
  }

  loadYouTubeAPI() {
    const tag = document.createElement('script');
    tag.src = IFRAME_API_URL;
    const firstScriptTag = document.getElementsByTagName('script')[0];
    firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

    // global callback
    window.onYouTubeIframeAPIReady = () => this.onYouTubeIframeAPIReady();
  }

  onYouTubeIframeAPIReady() {
    this.player = new YT.Player(this.playerElementId, {
      width: '100%',
      height: '100%',
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
        autoplay: 1,
        loop: 0
      }
    });
  }

  onPlayerReady(event) {
    this.playerReady = true;
    event.target.playVideo();
  }

  onPlayerStateChange(event) {
    console.log({ event, playerState: YT.PlayerState });

    if (event.data == YT.PlayerState.PLAYING) {
    }
  }

  stopVideo() {
    this.player.stopVideo();
    this.activateBuffering();
  }

  playVideo() {
    this.player.playVideo();
    this.deactivateBuffering();
  }

  activateBuffering = () => {};
  deactivateBuffering = () => {};

  loadChannels() {
    fetch(this.channelsUrl)
      .then((res) => res.json())
      .then((res) => {
        const channelsContainer = document.getElementById('channels');
        res.channels.forEach((channel) => {
          const button = document.createElement('button');
          button.textContent = channel.type;
          button.addEventListener('click', () =>
            this.loadChannelVideo(channel)
          );
          channelsContainer.appendChild(button);
        });
      });
  }

  loadChannelVideo(channel) {
    const videoId = channel.videos[0].url.split('=')[1];
    console.log('Loading video ID:', videoId);

    if (this.playerReady && this.player.loadVideoById) {
      this.player.loadVideoById(videoId);
    } else {
      console.error('Player not ready or invalid video ID:', videoId);
    }
  }

  addControlListeners() {
    document.querySelector('#control-pause').addEventListener('click', () => {
      if (this.playerReady) {
        this.stopVideo();
      }
    });
  }
}

const youtubePlayer = new YouTubePlayer('player');
