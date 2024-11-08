package db

import (
	"database/sql"
	"log"

	"github.com/ozencb/couchtube/helpers"
	jsonmodels "github.com/ozencb/couchtube/models/json"
	_ "modernc.org/sqlite"
)

func InitTables() {
	db, err := GetConnector()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createVideosTableQuery := `CREATE TABLE IF NOT EXISTS videos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"channel_id" INTEGER NOT NULL,
		"url" TEXT NOT NULL,
		"segment_start" INTEGER NOT NULL,
		"segment_end" INTEGER NOT NULL,
		UNIQUE(url, channel_id),
		FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE,
		CHECK (segment_end > segment_start)
	);`
	createChannelsTableQuery := `CREATE TABLE IF NOT EXISTS channels (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		UNIQUE(name)
	);`
	createIndexesQuery := `CREATE INDEX IF NOT EXISTS idx_videos_channel_id ON videos(channel_id);`

	_, err = db.Exec(createChannelsTableQuery + createVideosTableQuery + createIndexesQuery)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized and tables created successfully.")
}

func PopulateDatabase() {
	// parse the json file and insert the data into the database
	// ignore if there are channels already defined

	channels, err := helpers.LoadJSONFromFile[jsonmodels.ChannelsJson]("/channels.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", "file:couchtube.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check if anything already exists in channels
	var exists int
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM channels LIMIT 1);`).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists == 1 {
		log.Println("Data already exists in the database. Skipping db population")
		return
	}

	insertChannelQuery := `INSERT OR IGNORE INTO channels (name) VALUES (?)`
	insertVideoQuery := `INSERT OR IGNORE INTO videos (channel_id, url, segment_start, segment_end) VALUES (?, ?, ?, ?)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, channel := range channels.Channels {
		channelID, err := insertOrGetChannelID(tx, channel.Name, insertChannelQuery)
		if err != nil {
			log.Fatal(err)
		}

		for _, video := range channel.Videos {
			_, err = tx.Exec(insertVideoQuery, channelID, video.Url, video.SegmentStart, video.SegmentEnd)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.Fatal(err)
	}

	log.Println("Data inserted successfully.")
}

func insertOrGetChannelID(tx *sql.Tx, name, query string) (int64, error) {
	result, err := tx.Exec(query, name)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected > 0 {
		return result.LastInsertId()
	}

	var existingID int64
	getIDQuery := `SELECT id FROM channels WHERE name = ?`
	err = tx.QueryRow(getIDQuery, name).Scan(&existingID)
	if err != nil {
		return 0, err
	}

	return existingID, nil
}
