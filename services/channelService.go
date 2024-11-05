package services

import (
	"time"

	"github.com/ozencb/couchtube/db"
	dbmodels "github.com/ozencb/couchtube/models/db"
	repo "github.com/ozencb/couchtube/repositories"
)

func GetChannels() ([]dbmodels.Channel, error) {

	db, err := db.GetConnector()

	if err != nil {
		return nil, err
	}

	repo := repo.NewChannelRepository(db)

	channels, err := repo.GetChannels()

	if err != nil {
		return nil, err
	}

	return channels, nil
}

func GetCurrentVideoByChannelId(channelId int) (*dbmodels.Video, error) {

	db, err := db.GetConnector()

	if err != nil {
		return nil, err
	}

	repo := repo.NewVideoRepository(db)

	videos, err := repo.GetVideosByChannelID(channelId)

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
		if totalLengthInSeconds <= indexSeconds {
			if indexSeconds > 0 {
				video.SegmentStart = video.SegmentStart + indexSeconds // wind the video forward to the correct position
			}
			return &video, nil
		}
	}

	println("No video found to play")

	return &videos[0], nil
}
