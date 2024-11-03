package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func GetConnector() (*sql.DB, error) {
	dbFilePath := "couchtube.db"

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		file, err := os.Create(dbFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}
		file.Close()
		log.Printf("Created new SQLite database file at %s", dbFilePath)
	}

	dsn := fmt.Sprintf("file:%s?cache=shared", dbFilePath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// verify the connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
