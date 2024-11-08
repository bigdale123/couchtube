package services

import (
	"time"

	dbmodels "github.com/ozencb/couchtube/models/db"
	repo "github.com/ozencb/couchtube/repositories"
)

type ChannelService struct {
	ChannelRepo repo.ChannelRepository
	VideoRepo   repo.VideoRepository
}

func NewChannelService(channelRepo repo.ChannelRepository, videoRepo repo.VideoRepository) *ChannelService {
	return &ChannelService{
		ChannelRepo: channelRepo,
		VideoRepo:   videoRepo,
	}
}

func (s *ChannelService) GetChannels() ([]dbmodels.Channel, error) {
	channels, err := s.ChannelRepo.GetChannels()

	if err != nil {
		return nil, err
	}

	return channels, nil
}

func (s *ChannelService) GetCurrentVideoByChannelId(channelId int) (*dbmodels.Video, error) {

	videos, err := s.VideoRepo.GetVideosByChannelID(channelId)

	if err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		return nil, nil
	}

	// get total video length
	totalLengthInSeconds := 0
	for _, video := range videos {
		totalLengthInSeconds += video.SegmentEnd - video.SegmentStart
	}

	// get seconds elapsed since the beginning of the day
	secondsElapsed := time.Now().Hour()*3600 + time.Now().Minute()*60 + time.Now().Second()

	indexSeconds := secondsElapsed % totalLengthInSeconds

	// find the video to be played and update the segment start time
	for _, video := range videos {
		totalLengthInSeconds -= video.SegmentEnd - video.SegmentStart
		if totalLengthInSeconds <= indexSeconds && indexSeconds < (video.SegmentEnd-video.SegmentStart) {
			if indexSeconds > 0 {
				video.SegmentStart = video.SegmentStart + indexSeconds // wind the video forward to the correct position
			}
			return &video, nil
		}
	}

	println("No video found to play")

	return &videos[0], nil
}

func (s *ChannelService) GetNextVideo(channelId int, videoId int) *dbmodels.Video {
	video, err := s.VideoRepo.GetNextVideo(channelId, videoId)
	if err != nil {
		return nil
	}

	return video
}
