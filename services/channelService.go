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

	totalLength := int64(0)
	for _, video := range videos {
		totalLength += int64(video.SegmentEnd - video.SegmentStart)
	}

	currentPoint := time.Now().UTC().Unix() % totalLength
	videoIndex := -1

	for i := range videos {
		video := &videos[i]
		segmentLength := int64(video.SegmentEnd - video.SegmentStart)

		if currentPoint < segmentLength {
			videoIndex = i
			video.SegmentStart += int(currentPoint) // Adjust start to match the current second
			break
		}
		currentPoint -= segmentLength
	}

	if videoIndex == -1 {
		return &videos[0], nil
	}

	return &videos[videoIndex], nil
}

func (s *ChannelService) GetNextVideo(channelId int, videoId int) *dbmodels.Video {
	video, err := s.VideoRepo.GetNextVideo(channelId, videoId)
	if err != nil {
		return nil
	}

	return video
}
