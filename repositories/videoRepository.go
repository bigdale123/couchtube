package repo

import (
	"database/sql"

	dbmodels "github.com/ozencb/couchtube/models/db"
)

type VideoRepository interface {
	GetVideosByChannelID(channelID int) ([]dbmodels.Video, error)
}

type videoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) VideoRepository {
	return &videoRepository{db: db}
}

func (r *videoRepository) GetVideosByChannelID(channelID int) ([]dbmodels.Video, error) {
	rows, err := r.db.Query(`
        SELECT id, channel_id, url, segment_start, segment_end
        FROM videos
        WHERE channel_id = ?
    `, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []dbmodels.Video
	for rows.Next() {
		var video dbmodels.Video
		if err := rows.Scan(&video.ID, &video.ChannelID, &video.URL, &video.SegmentStart, &video.SegmentEnd); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}
