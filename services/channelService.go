package services

import (
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

	return &videos[0], nil
}
