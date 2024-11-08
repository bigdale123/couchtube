package services

import (
	"encoding/json"
	"net/http"

	jsonmodels "github.com/ozencb/couchtube/models/json"
	repo "github.com/ozencb/couchtube/repositories"
)

type SubmitListService struct {
	ChannelRepo repo.ChannelRepository
	VideoRepo   repo.VideoRepository
}

func NewSubmitListService(channelRepo repo.ChannelRepository, videoRepo repo.VideoRepository) *SubmitListService {
	return &SubmitListService{
		ChannelRepo: channelRepo,
		VideoRepo:   videoRepo,
	}
}

func (s *SubmitListService) SubmitList(list jsonmodels.SubmitListRequestJson) (bool, error) {
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

	s.ChannelRepo.PurgeChannels()

	s.VideoRepo.PurgeVideos()

	for _, channel := range videoList.Channels {
		println("Inserting channel:", channel.Name)
		channelId, err := s.ChannelRepo.SaveChannel(channel.Name)

		if err != nil {
			return false, err
		}
		for _, video := range channel.Videos {
			println("Inserting video:", channelId, video.Url, video.SegmentStart, video.SegmentEnd)
			s.VideoRepo.SaveVideo(channelId, video.Url, video.SegmentStart, video.SegmentEnd)
		}
	}

	return true, nil
}
