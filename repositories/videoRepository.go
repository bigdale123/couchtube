package repo

import (
	"database/sql"

	"github.com/ozencb/couchtube/db"
	dbmodels "github.com/ozencb/couchtube/models/db"
)

type VideoRepository interface {
	GetVideosByChannelID(channelID int) ([]dbmodels.Video, error)
	GetNextVideo(videoID int, channelID int) (*dbmodels.Video, error)
}

type videoRepository struct {
	db *sql.DB
}

func NewVideoRepository() VideoRepository {
	db, err := db.GetConnector()

	if err != nil {
		panic(err)
	}

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

func (r *videoRepository) GetNextVideo(videoID int, channelID int) (*dbmodels.Video, error) {
	row := r.db.QueryRow(`
		SELECT id, channel_id, url, segment_start, segment_end
		FROM videos
		WHERE channel_id = ? AND id > ?
		ORDER BY id ASC
		LIMIT 1
	`, channelID, videoID)

	var video dbmodels.Video
	err := row.Scan(&video.ID, &video.ChannelID, &video.URL, &video.SegmentStart, &video.SegmentEnd)
	if err == sql.ErrNoRows {
		// If no next video is found, get the first video instead
		row = r.db.QueryRow(`
			SELECT id, channel_id, url, segment_start, segment_end
			FROM videos
			WHERE channel_id = ?
			ORDER BY id ASC
			LIMIT 1
		`, channelID)

		err = row.Scan(&video.ID, &video.ChannelID, &video.URL, &video.SegmentStart, &video.SegmentEnd)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &video, nil
}
