package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	port        string
	dbFilePath  string
	initialized bool
	readonly    bool
	once        sync.Once
)

func LoadEnv() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, relying on system environment variables")
		}

		port = getEnv("PORT", "8081")
		dbFilePath = getEnv("DATABASE_FILE", "couchtube.db")
		readonly = getEnv("READONLY_MODE", "false") == "true"
		initialized = true
	})
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetPort() string {
	if !initialized {
		LoadEnv()
	}
	return port
}

func GetDBFilePath() string {
	if !initialized {
		LoadEnv()
	}
	return dbFilePath
}

func GetReadonlyMode() bool {
	if !initialized {
		LoadEnv()
	}
	return readonly
}
