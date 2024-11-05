package repo

import (
	"database/sql"

	"github.com/ozencb/couchtube/db"
	dbmodels "github.com/ozencb/couchtube/models/db"
)

type ChannelRepository interface {
	GetChannels() ([]dbmodels.Channel, error)
}

type channelRepository struct {
	db *sql.DB
}

func NewChannelRepository() ChannelRepository {
	db, err := db.GetConnector()

	if err != nil {
		panic(err)
	}
	return &channelRepository{db: db}
}

func (r *channelRepository) GetChannels() ([]dbmodels.Channel, error) {
	rows, err := r.db.Query("SELECT id, name FROM channels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []dbmodels.Channel
	for rows.Next() {
		var channel dbmodels.Channel
		if err := rows.Scan(&channel.ID, &channel.Name); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return channels, nil
}
