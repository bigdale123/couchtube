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

	// Fetch the video list from the provided URL
	response, err := http.Get(videoListUrl)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	var videoList jsonmodels.ChannelsJson
	err = json.NewDecoder(response.Body).Decode(&videoList)
	if err != nil {
		return false, err
	}

	if len(videoList.Channels) == 0 {
		return false, nil
	}

	// Begin a transaction
	tx, err := s.ChannelRepo.BeginTx()
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback() // Roll back if there's an error
		} else {
			tx.Commit() // Commit if everything succeeds
		}
	}()

	// Purge existing channels and videos
	if err = s.ChannelRepo.PurgeChannels(tx); err != nil {
		return false, err
	}
	if err = s.VideoRepo.PurgeVideos(tx); err != nil {
		return false, err
	}

	// Insert new channels and videos
	for _, channel := range videoList.Channels {
		channelID, err := s.ChannelRepo.SaveChannel(tx, channel.Name)
		if err != nil {
			return false, err
		}
		for _, video := range channel.Videos {
			err = s.VideoRepo.SaveVideo(tx, channelID, video.Url, video.SegmentStart, video.SegmentEnd)
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
