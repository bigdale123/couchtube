package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	dbInstance *sql.DB
	once       sync.Once
	closeOnce  sync.Once
)

func GetDbConnection() (*sql.DB, error) {
	var err error
	once.Do(func() {
		dbFilePath := "couchtube.db"

		// Create the database file if it doesn't exist
		if _, err = os.Stat(dbFilePath); os.IsNotExist(err) {
			var file *os.File
			file, err = os.Create(dbFilePath)
			if err != nil {
				log.Fatalf("Failed to create database file: %v", err)
			}
			file.Close()
			log.Printf("Created new SQLite database file at %s", dbFilePath)
		}

		// Open the database
		dsn := fmt.Sprintf("file:%s?cache=shared", dbFilePath)
		dbInstance, err = sql.Open("sqlite", dsn)
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}

		// Enable foreign key constraints
		if _, err = dbInstance.Exec("PRAGMA foreign_keys = ON;"); err != nil {
			log.Fatalf("Failed to enable foreign keys: %v", err)
		}

		// Test the database connection
		if err = dbInstance.Ping(); err != nil {
			dbInstance.Close()
			log.Fatalf("Failed to ping database: %v", err)
		}
	})

	if err != nil {
		log.Printf("Failed to get database connector: %v", err)
		return nil, err
	}

	return dbInstance, err
}

func CloseConnector() {
	closeOnce.Do(func() {
		if dbInstance != nil {
			if err := dbInstance.Close(); err != nil {
				log.Printf("Error closing the database connection: %v", err)
			} else {
				log.Println("Database connection closed successfully.")
			}
			dbInstance = nil // Ensure we don't try to close it again
		}
	})
}
