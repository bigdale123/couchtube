package services

import (
	"encoding/json"
	"net/http"

	jsonmodels "github.com/ozencb/couchtube/models/json"
	repo "github.com/ozencb/couchtube/repositories"
)

func SubmitList(list jsonmodels.SubmitListRequestJson) (bool, error) {
	videoListUrl := list.VideoListUrl

	if videoListUrl == "" {
		return false, nil
	}

	response, err := http.Get(videoListUrl)
	if err != nil {
		return false, err
	}

	var videoList jsonmodels.ChannelsJson
	err = json.NewDecoder(response.Body).Decode(&videoList)
	if err != nil {
		return false, err
	}

	if len(videoList.Channels) == 0 {
		return false, nil
	}

	channelRepo := repo.NewChannelRepository()
	videoRepo := repo.NewVideoRepository()

	channelRepo.PurgeChannels()
	videoRepo.PurgeVideos()

	for _, channel := range videoList.Channels {
		channelId, err := channelRepo.SaveChannel(channel.Name)

		if err != nil {
			return false, err
		}
		for _, video := range channel.Videos {
			println(channelId, video.Url, video.SegmentStart, video.SegmentEnd)
			videoRepo.SaveVideo(channelId, video.Url, video.SegmentStart, video.SegmentEnd)
		}
	}

	return true, nil
}
