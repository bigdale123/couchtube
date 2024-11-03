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

	createVideosTable := `CREATE TABLE IF NOT EXISTS videos (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"channel_id" INTEGER NOT NULL,
		"url" TEXT NOT NULL,
		"segment_start" INTEGER NOT NULL,
		"segment_end" INTEGER NOT NULL,
		UNIQUE(url, channel_id),
		FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE
	);`

	createChannelsTable := `CREATE TABLE IF NOT EXISTS channels (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT
	);`

	createIndexes := `CREATE INDEX IF NOT EXISTS idx_videos_channel_id ON videos(channel_id);`

	_, err = db.Exec(createChannelsTable + createVideosTable + createIndexes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized and table created successfully.")
}

func PopulateDatabase() {
	// parse the json file and insert the data into the database

	channels, err := helpers.LoadJSONFromFile[jsonmodels.ChannelsJson]("/channels.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", "file:couchtube.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	insertChannel := `INSERT INTO channels (name) VALUES (?) RETURNING id`
	insertVideo := `INSERT INTO videos (channel_id, url, segment_start, segment_end) VALUES (?, ?, ?, ?)`

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(insertChannel)
	if err != nil {
		log.Fatal(err)
	}

	for _, channel := range channels.Channels {
		result, err := stmt.Exec(channel.Name)
		if err != nil {
			log.Fatal(err)
		}
		channelId, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		if err != nil {
			log.Fatal(err)
		}

		for _, video := range channel.Videos {
			_, err = tx.Exec(insertVideo, channelId, video.Url, video.SegmentStart, video.SegmentEnd)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	stmt.Close()

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Data inserted successfully.")

}
